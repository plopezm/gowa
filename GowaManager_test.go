package gowa

import (
	"testing"
)

type User struct{
	Email		string	`gorm:"primary_key" gowa:"pk"`
	Password	string
	Role		string
}

type Company struct{
	User		User	`gorm:"ForeignKey:user_email;AssociationForeignKey:email" gowa:"ignore" json:"omitempty"`
	UserEmail	string	`gorm:"primary_key" gowa:"pk;fk_table:User;fk_col:Email"`
	Name		string
	Cif		string	`gorm:"unique"`
	Sector		string
	Tlfn		string
	Country		string
	City		string
	Address		string
}

type Driver struct {
	User		User	`gorm:"ForeignKey:user_email;AssociationForeignKey:email" gowa:"ignore" json:"omitempty"`
	UserEmail	string	`gorm:"primary_key" gowa:"pk;fk_table:User;fk_col:Email"`
	Name		string	`gorm:"unique"`
	Lastname	string
	Cname		string
	Cphone		string
	Ccity		string
	Ccountry	string
	Ccp		string
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
			if !(value.Name == "User" && value.Ctype == "User"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Name == "UserEmail" && value.Ctype == "string" && value.Pk && value.Fktab == "User" && value.Fkcol == "Email"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.Name == "Name" && value.Ctype == "string" && !value.Pk && value.Fktab == "" && value.Fkcol == ""){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

}