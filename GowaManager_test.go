package gowa

import (
	"testing"
)

type User struct{
	Email		string	`goedb:"pk"`
	Password	string
	Role		string
}

type Company struct{
	UserEmail	string	`goedb:"pk"`
	Name		string
	Cif		string	`goedb:"unique"`
}

type Driver struct {
	UserEmail	string	`goedb:"pk,fk=User(Email)"`
	Name		string	`goedb:"unique"`
}

type Vehicle struct {
	Owner 		string  `goedb:"pk,fk=Driver(UserEmail)"`
	Model		string
	Plate		string	`goedb:"unique"`
}

var am *GowaManager

func TestGowaManager_AddModel(t *testing.T) {
	am = NewGowa("sqlite3", "test.db", 10)

	am.Open()
	defer am.Close()

	err := am.AddModel(&User{})
	if err != nil {
		t.Error(err)
	}

	err = am.AddModel(&Company{})
	if err != nil {
		t.Error(err)
	}

	err = am.AddModel(&Driver{})
	if err != nil {
		t.Error(err)
	}

	err = am.AddModel(&Vehicle{})
	if err != nil {
		t.Error(err)
	}
}


func TestGowaManager_RemoveModel(t *testing.T) {

	am.Open()
	defer am.Close()

	err := am.RemoveModel(&Company{})
	if err != nil {
		t.Error(err)
	}
	err = am.RemoveModel(&Vehicle{})
	if err != nil {
		t.Error(err)
	}
	err = am.RemoveModel(&Driver{})
	if err != nil {
		t.Error(err)
	}
	err = am.RemoveModel(&User{})
	if err != nil {
		t.Error(err)
	}

	_, err = am.db.Model(&User{})
	if err == nil {
		t.Error("Model exists")
	}
}

/*
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


	gowaTable.Model, gowaTable.Title, gowaTable.Columns = parseModel(&Vehicle{})

	if gowaTable.Title != "Driver"{
		t.Error("Error getting model name")
	}

	for key, value := range gowaTable.Columns {
		switch key{
		case 0:
			if !(value.Name == "Owner" && value.Ctype == "string" && value.Pk && value.Fktab == "Driver" && value.Fkcol == "UserEmail") {
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

}*/
