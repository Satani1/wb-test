package app

import "wb-test/pkg/db"

type Application struct {
	Addr   string
	DB     *db.Repository
	Secret string
}
