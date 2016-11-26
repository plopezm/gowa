package gowa

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/plopezm/goServerUtils"
	"reflect"
	"strings"
)

var GM *GowaManager

type GowaManager struct {
	db          *gorm.DB
	adminTables map[string]GowaTable
	dbType      string
	dbPath      string
	PageSize    uint32
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
		panic( txt )
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

func isComposed(v interface{}) (bool){

	switch t:= v.(type){

	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64, float32, float64:
		return false;
	case string:
		return false;
	default:
		_ = t;
		return true;
	}
}

func manageTag(gowacol *GowaColumn,tag string){
	if strings.Contains(tag, ";") {
		attrs := strings.Split(tag, ";")
		for i:=0;i<len(attrs);i++ {
			if strings.Contains(attrs[i], ":") {
				fieldVal := strings.Split(attrs[i], ":")
				switch(fieldVal[0]){
				case "fk_table":
					gowacol.Fktab = fieldVal[1]
				case "fk_col":
					gowacol.Fkcol = fieldVal[1]
				}
				continue
			}
			switch attrs[i] {
			case "pk":
				gowacol.Pk = true
			}
		}
	}

	switch tag {
	case "pk":
		gowacol.Pk = true
	}
}

func parseModel(model interface{}) (reflect.Type, string, []GowaColumn){
	typ := reflect.TypeOf(model)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	columnSlice := make([]GowaColumn, typ.NumField())

	for i:=0;i<typ.NumField();i++ {
		gowacol := GowaColumn{}
		gowacol.Name = typ.Field(i).Name
		gowacol.Ctype = typ.Field(i).Type.Name()

		if val, ok := typ.Field(i).Tag.Lookup("gowa"); ok {
			manageTag(&gowacol, val)
		}

		columnSlice[i] = gowacol
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

func (am *GowaManager) getRoutes() goServerUtils.Routes {
	routes := goServerUtils.Routes{
		goServerUtils.Route{
			"ValidateSession",
			"GET",
			"/gowa/api/validate",
			ValidateSession,
		},
		goServerUtils.Route{
			"LoginAdmin",
			"GET",
			"/gowa/api/login",
			Login,
		},
		goServerUtils.Route{
			"CreateAdmin",
			"PUT",
			"/gowa/api/register",
			CreateUser,
		},
		goServerUtils.Route{
			"GetTablesStruct",
			"GET",
			"/gowa/api/rest/tables",
			GetTablesStruct,
		},
		goServerUtils.Route{
			"GetTableStruct",
			"GET",
			"/gowa/api/rest/tables/struct/{table}",
			GetTableStruct,
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


