package gowa

import (
	"testing"
)

type User struct{
	Email		string	`goedb:"pk" gowa:"pk"`
	Password	string
	Role		string
}

type Company struct{
	UserEmail	string	`goedb:"pk" gowa:"pk;fk_table:User;fk_col:Email"`
	Name		string
	Cif		string	`goedb:"unique"`
}

type Driver struct {
	UserEmail	string	`goedb:"pk,fk=User(Email)" gowa:"pk;fk_table:User;fk_col:Email"`
	Name		string	`goedb:"unique"`
}

type Vehicle struct {
	Owner 		string  `goedb:"pk,fk=Driver(UserEmail)"`
	Model		string
	Plate		string	`goedb:"unique"`
}


func TestGowaParseModel(t *testing.T){
	var gowaTable GowaTable

	gowaTable.Model, gowaTable.Title, gowaTable.Columns = parseModel(&User{})

	if gowaTable.Title != "User"{
		t.Error("Error getting model name")
	}

	for key, value := range gowaTable.Columns {
		switch key{
		case 0:
			if !(value.Name == "Email" && value.Ctype == "string" && value.Pk){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Name == "Password" && value.Ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.Name == "Role" && value.Ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

	gowaTable.Model, gowaTable.Title, gowaTable.Columns = parseModel(&Company{})

	if gowaTable.Title != "Company"{
		t.Error("Error getting model name")
	}

	for key, value := range gowaTable.Columns {
		switch key{
		case 0:
			if !(value.Name == "UserEmail" && value.Ctype == "string" && value.Pk && value.Fktab == "User" && value.Fkcol == "Email"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Name == "Name" && value.Ctype == "string" && !value.Pk && value.Fktab == "" && value.Fkcol == ""){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}


	gowaTable.Model, gowaTable.Title, gowaTable.Columns = parseModel(&Driver{})

	if gowaTable.Title != "Driver"{
		t.Error("Error getting model name")
	}

	for key, value := range gowaTable.Columns {
		switch key{
		case 0:
			if !(value.Name == "UserEmail" && value.Ctype == "string" && value.Pk && value.Fktab == "User" && value.Fkcol == "Email") {
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Name == "Name" && value.Ctype == "string" && !value.Pk && value.Fktab == "" && value.Fkcol == "") {
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

}