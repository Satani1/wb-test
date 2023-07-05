package models

type Task struct {
	ID     int
	Item   string `fake:"{randomstring:[sand,soil,tv,ATM,kettlebell,vintage clock,bed,sofa,chair]}"`
	Weight int    `fake:"{number:10,80}"`
}
