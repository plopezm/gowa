package gowa

import (
	"reflect"
	"github.com/plopezm/goedb"
)

type GowaColumn struct {
	Name 		string			`json:"name"`
	Ctype 		string			`json:"ctype"`
	Pk		bool			`json:"pk"`
	Fktab		string			`json:"fktab"`
	Fkcol		string			`json:"fkcol"`
}

type GowaTable struct {
	Title		string			`json:"title"`
	Columns		[]goedb.GoedbColumn	`json:"columns"`
	Rows 		interface{}  		`json:"rows"`
	Page		uint64			`json:"page"`

	Model		reflect.Type		`json:"-"`
}
