package gowa

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/plopezm/goServerUtils"
	"reflect"
)

var GM *GowaManager

type GowaManager struct {
	DB          *gorm.DB
	AdminTables map[string]GowaTable
	DbType      string
	DbPath      string
	pageSize    uint32
}


func (am *GowaManager) Init(dbtype string, dbpath string, pageSize uint32) error{
	var err error;
	am.DB, err = gorm.Open(dbtype, dbpath);
	if err != nil {
		panic(err)
	}
	defer am.DB.Close();

	if err != nil {
		return err;
	}
	am.DbPath = dbpath
	am.DbType = dbtype
	am.AdminTables = make(map[string]GowaTable)
	am.pageSize = pageSize
	return nil
}

func (am *GowaManager) GetSession() (*gorm.DB, error){
	var err error;

	if am.DB.DB().Ping() != nil{
		am.DB, err = gorm.Open(am.DbType, am.DbPath);
		return am.DB, err;
	}
	return am.DB, nil;
}

func (am *GowaManager) End(){
	am.DB.Close();
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

	am.AdminTables[gowaTable.Title] = gowaTable
}

func (am *GowaManager) RemoveModel(table_name string){
	delete(am.AdminTables, table_name)
	//delete(am.AdminModels, table_name)
}

func (am *GowaManager) GetRoutes() goServerUtils.Routes {
	routes := goServerUtils.Routes{
		goServerUtils.Route{
			"GetTablesStruct",
			"GET",
			"/gowa/api/rest/tables",
			GetTablesStruct,
		},
		goServerUtils.Route{
			"GetTable",
			"GET",
			"/gowa/api/rest/tables/show/{table}",
			GetTable,
		},
		goServerUtils.Route{
			"ShowTableRows",
			"GET",
			"/gowa/api/rest/tables/show/rows/{table}",
			ShowTableRows,
		},
		goServerUtils.Route{
			"AddTableRow",
			"PUT",
			"/gowa/api/rest/tables/add/row/{table}",
			AddTableRow,
		},
		goServerUtils.Route{
			"RemoveTableRow",
			"DELETE",
			"/gowa/api/rest/tables/remove/row/{table}",
			RemoveTableRow,
		},
	};
	return routes
}


