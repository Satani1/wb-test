package models

type Loader struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	MaxWeight int    `json:"max-weight" fake:"{number:5,30}"`
	Drunk     bool   `json:"drunk" fake:"{bool}"`
	Fatigue   int    `json:"fatigue" fake:"{number:0,100}"`
	Salary    int    `json:"salary" fake:"{number:10000,30000}"`
}
