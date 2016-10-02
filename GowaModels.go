package gowa


type GowaTableRow struct{
	values		[]interface{}

}

type GowaTable struct {
	Title		string		`json:"title"`
	Columns		[]string	`json:"columns"`
	Rows 		[]GowaTableRow  `json:"rows"`
}