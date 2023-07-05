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
	stmt := `INSERT INTO loadertable (username, password) values (?,?)`

	result, err := r.db.Exec(stmt, loader.Username, loader.Password)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// GetLoader return Loader struct by a username from table
func (r *Repository) GetLoader(username string) (*models.Loader, error) {
	stmt := `select username,password from loadertable where username = ?`

	row := r.db.QueryRow(stmt, username)
	var loader models.Loader
	if err := row.Scan(&loader.Username, &loader.Password); err != nil {
		return nil, err
	}

	return &loader, nil
}

// InsertCustomer insert a row into Customer table
func (r *Repository) InsertCustomer(customer models.Customer) (int, error) {
	stmt := `INSERT INTO customertable (username, password) values (?,?)`

	result, err := r.db.Exec(stmt, customer.Username, customer.Password)
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
	stmt := `select username,password from customertable where username = ?`

	row := r.db.QueryRow(stmt, username)
	var customer models.Customer
	if err := row.Scan(&customer.Username, &customer.Password); err != nil {
		return nil, err
	}

	return &customer, nil
}
