package db

import "wb-test/pkg/models"

// InsertCustomer insert a row into Customer table
func (r *Repository) InsertCustomer(customer models.Customer) error {
	stmt := `INSERT INTO customertable (username, password,capital) values ($1,$2,$3)`

	_, err := r.db.Exec(stmt, customer.Username, customer.Password, customer.StartCapital)
	if err != nil {
		return err
	}

	return nil
}

// GetCustomer return Customer struct by a username from table
func (r *Repository) GetCustomer(username string) (*models.Customer, error) {
	stmt := `select ID,capital,username,password from customertable where username = $1`

	row := r.db.QueryRow(stmt, username)
	var customer models.Customer
	if err := row.Scan(&customer.ID, &customer.StartCapital, &customer.Username, &customer.Password); err != nil {
		return nil, err
	}

	return &customer, nil
}

// UpdateCustomer updates capital data in customer row
func (r *Repository) UpdateCustomer(id int, newCapital int) error {
	stmt := `update customertable set capital = $1 where  ID = $2`

	_, err := r.db.Exec(stmt, newCapital, id)
	if err != nil {
		return err
	}
	return nil
}
