package models

type Loader struct {
	ID        int
	Username  string
	Password  string
	MaxWeight int  `fake:"{number:5,30}"`
	Drunk     bool `fake:"{bool}"`
	Fatigue   int  `fake:"{number:0,100}"`
	Salary    int  `fake:"{number:10000,30000}"`
}
