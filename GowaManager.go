package gowa

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/plopezm/goServerUtils"
	"reflect"
	"fmt"
	"net/http"
)

var GM *GowaManager

type GowaManager struct {
	db          *gorm.DB
	adminTables map[string]GowaTable
	dbType      string
	dbPath      string
	PageSize    uint32
}

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}


func (am *GowaManager) init(dbtype string, dbpath string, pageSize uint32) error{
	var err error;
	am.db, err = gorm.Open(dbtype, dbpath)
	if err != nil {
		panic(err)
	}
	defer am.db.Close();

	if err != nil {
		return err;
	}
	am.dbPath = dbpath
	am.dbType = dbtype
	am.adminTables = make(map[string]GowaTable)
	am.PageSize = pageSize

	if err := am.db.AutoMigrate(GowaUser{}).Error; err != nil {
		txt := "AutoMigrate Job table failed"
		panic( fmt.Sprintf( "%s: %s", txt, err ) )
	}

	return nil
}

func (am *GowaManager) getSession() (*gorm.DB, error){
	var err error;

	if am.db.DB().Ping() != nil{
		am.db, err = gorm.Open(am.dbType, am.dbPath);
		return am.db, err;
	}
	return am.db, nil;
}

func (am *GowaManager) End(){
	am.db.Close();
}

func parseModel(model interface{}) (reflect.Type, string, []string){
	typ := reflect.TypeOf(model)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	columnSlice := make([]string, typ.NumField())

	for i:=0;i<typ.NumField();i++ {
		columnSlice[i] = typ.Field(i).Name
	}

	return typ, typ.Name(), columnSlice
}

func (am *GowaManager) AddModel(model interface{}){
	var gowaTable GowaTable

	gowaTable.Model, gowaTable.Title, gowaTable.Columns = parseModel(model)

	am.adminTables[gowaTable.Title] = gowaTable
}

func (am *GowaManager) RemoveModel(table_name string){
	delete(am.adminTables, table_name)
}

func (am *GowaManager) getRoutes() []Route {
	routes := []Route{
		Route{
			"ValidateSession",
			"GET",
			"/gowa/api/validate",
			ValidateSession,
		},
		Route{
			"LoginAdmin",
			"GET",
			"/gowa/api/login",
			Login,
		},
		Route{
			"CreateAdmin",
			"PUT",
			"/gowa/api/register",
			CreateUser,
		},
		Route{
			"GetTablesStruct",
			"GET",
			"/gowa/api/rest/tables",
			GetTablesStruct,
		},
		Route{
			"GetTable",
			"GET",
			"/gowa/api/rest/tables/show/{table}",
			GetTable,
		},
		Route{
			"ShowTableRows",
			"GET",
			"/gowa/api/rest/tables/show/rows/{table}",
			ShowTableRows,
		},
		Route{
			"AddTableRow",
			"PUT",
			"/gowa/api/rest/tables/add/row/{table}",
			AddTableRow,
		},
		Route{
			"RemoveTableRow",
			"DELETE",
			"/gowa/api/rest/tables/remove/row/{table}",
			RemoveTableRow,
		},
	};
	return routes
}


