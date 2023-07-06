package db

import "wb-test/pkg/models"

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

// UpdateCustomer updates capital data in customer row
func (r *Repository) UpdateCustomer(id int, newCapital int) error {
	stmt := `update customertable set capital = ? where  ID = ?`

	_, err := r.db.Exec(stmt, newCapital, id)
	if err != nil {
		return err
	}
	return nil
}
