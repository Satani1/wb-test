package db

import "wb-test/pkg/models"

// InsertTask insert a row into tasks table
func (r *Repository) InsertTask(task models.Task) error {
	stmt := `INSERT INTO tasks (name, weight) values ($1,$2)`

	if _, err := r.db.Exec(stmt, task.Item, task.Weight); err != nil {
		return err
	}

	return nil
}

// GetTask return a one task by ID
func (r *Repository) GetTask(id int) (*models.Task, error) {
	stmt := `select ID, name, weight from tasks where ID = $1`
	var task models.Task

	row := r.db.QueryRow(stmt, id)

	if err := row.Scan(&task.ID, &task.Item, &task.Weight); err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTaskAvailable return all task that available
func (r *Repository) GetTaskAvailable() ([]models.Task, error) {
	stmt := `select ID, name, weight from tasks where done = 0`
	var tasks []models.Task

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.Item, &task.Weight)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTaskCompleted return tasks that completed by loader (need loader's ID)
func (r *Repository) GetTaskCompleted(id int) ([]models.Task, error) {
	stmt := `select distinct tasks.ID, tasks.name, tasks.weight from donetasks, tasks inner join donetasks d on tasks.ID = d.task_ID where d.loader_ID = $1`
	var tasks []models.Task

	rows, err := r.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.Item, &task.Weight)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask update done from 0 to 1
// and INSERT into donetasks table that task
func (r *Repository) UpdateTask(taskID int, loaderID int) error {
	stmt1 := `update tasks set done = $1 where  ID = $2`
	stmt2 := `insert into donetasks (task_ID, loader_ID) values ($1, $2) `

	_, err := r.db.Exec(stmt1, 1, taskID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(stmt2, taskID, loaderID)
	if err != nil {
		return err
	}
	return nil
}
