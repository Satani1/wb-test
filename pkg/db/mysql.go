package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"wb-test/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewDB(url string) (*Repository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{
		db,
	}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

// InsertLoader insert a row into Loader table
func (r *Repository) InsertLoader(loader models.Loader) (int, error) {
	stmt := `INSERT INTO loadertable (username, password, MaxWeight, Drunk, Fatigue, Salary) values (?,?,?,?,?,?)`

	result, err := r.db.Exec(stmt, loader.Username, loader.Password, loader.MaxWeight, loader.Drunk, loader.Fatigue, loader.Salary)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// GetLoaderByName return Loader struct by a username from table
func (r *Repository) GetLoaderByName(username string) (*models.Loader, error) {
	stmt := `select ID, MaxWeight, Drunk, Fatigue, Salary, username, password from loadertable where username = ?`

	row := r.db.QueryRow(stmt, username)
	var loader models.Loader
	if err := row.Scan(&loader.ID, &loader.MaxWeight, &loader.Drunk, &loader.Fatigue, &loader.Salary, &loader.Username, &loader.Password); err != nil {
		return nil, err
	}

	return &loader, nil
}

// GetLoaderByID return Loader struct by an ID from table
func (r *Repository) GetLoaderByID(id int) (*models.Loader, error) {
	stmt := `select ID, MaxWeight, Drunk, Fatigue, Salary, username from loadertable where id = ?`

	row := r.db.QueryRow(stmt, id)

	var loader models.Loader
	if err := row.Scan(&loader.ID, &loader.MaxWeight, &loader.Drunk, &loader.Fatigue, &loader.Salary, &loader.Username); err != nil {
		return nil, err
	}

	return &loader, nil
}

// GetLoaders return all exist loaders
func (r *Repository) GetLoaders() ([]models.Loader, error) {
	stmt := `select ID, username, MaxWeight, Drunk, Fatigue, Salary from loadertable`
	var loaders []models.Loader

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var loader models.Loader

		err := rows.Scan(&loader.ID, &loader.Username, &loader.MaxWeight, &loader.Drunk, &loader.Fatigue, &loader.Salary)
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}
	return loaders, nil
}

// InsertCustomer insert a row into Customer table
func (r *Repository) InsertCustomer(customer models.Customer) (int, error) {
	stmt := `INSERT INTO customertable (username, password,capital) values (?,?,?)`

	result, err := r.db.Exec(stmt, customer.Username, customer.Password, customer.StartCapital)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// GetCustomer return Customer struct by a username from table
func (r *Repository) GetCustomer(username string) (*models.Customer, error) {
	stmt := `select ID,capital,username,password from customertable where username = ?`

	row := r.db.QueryRow(stmt, username)
	var customer models.Customer
	if err := row.Scan(&customer.ID, &customer.StartCapital, &customer.Username, &customer.Password); err != nil {
		return nil, err
	}

	return &customer, nil
}

// InsertTask insert a row into tasks table
func (r *Repository) InsertTask(task models.Task) error {
	stmt := `INSERT INTO tasks (name, weight) values (?,?)`

	if _, err := r.db.Exec(stmt, task.Item, task.Weight); err != nil {
		return err
	}

	return nil
}

// GetTask return a one task by ID
func (r *Repository) GetTask(id int) (*models.Task, error) {
	stmt := `select ID, name, weight from tasks where ID = ?`
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
	stmt := `select ID, name, weight from tasks where done = ?`
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
