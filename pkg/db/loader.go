package db

import "wb-test/pkg/models"

// InsertLoader insert a row into Loader table
func (r *Repository) InsertLoader(loader models.Loader) error {
	stmt := `INSERT INTO loadertable (username, password, MaxWeight, Drunk, Fatigue, Salary) values ($1,$2,$3,$4,$5,$6)`

	_, err := r.db.Exec(stmt, loader.Username, loader.Password, loader.MaxWeight, loader.Drunk, loader.Fatigue, loader.Salary)
	if err != nil {
		return err
	}

	return nil
}

// GetLoaderByName return Loader struct by a username from table
func (r *Repository) GetLoaderByName(username string) (*models.Loader, error) {
	stmt := `select ID, MaxWeight, Drunk, Fatigue, Salary, username, password from loadertable where username = $1`

	row := r.db.QueryRow(stmt, username)
	var loader models.Loader
	if err := row.Scan(&loader.ID, &loader.MaxWeight, &loader.Drunk, &loader.Fatigue, &loader.Salary, &loader.Username, &loader.Password); err != nil {
		return nil, err
	}

	return &loader, nil
}

// GetLoaderByID return Loader struct by an ID from table
func (r *Repository) GetLoaderByID(id int) (*models.Loader, error) {
	stmt := `select ID, MaxWeight, Drunk, Fatigue, Salary, username from loadertable where id = $1`

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

func (r *Repository) UpdateLoader(fatigue, id int) error {
	stmt := `update loadertable set Fatigue = $1 where ID = $2`

	if _, err := r.db.Exec(stmt, fatigue, id); err != nil {
		return err
	}
	return nil
}
