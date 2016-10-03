package gowa

type GowaTable struct {
	Title		string		`json:"title"`
	Columns		[]string	`json:"columns"`
	Rows 		interface{}  `json:"rows"`

	Page		uint64		`json:"page"`
	PageSize	uint8		`json:"pageSize"`
}