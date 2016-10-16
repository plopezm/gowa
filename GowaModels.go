package gowa

import "reflect"

type GowaTable struct {
	Title		string		`json:"title"`
	Columns		[]string	`json:"columns"`
	Rows 		interface{}  	`json:"rows"`

	Page		uint64		`json:"page"`

	Model		reflect.Type	`json:"-"`
}