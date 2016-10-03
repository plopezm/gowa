package gowa

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/plopezm/goServerUtils"
)

var GM *GowaManager

type GowaManager struct {
	DB          *gorm.DB
	AdminTables map[string]GowaTable
	AdminModels map[string]interface{}
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
	am.AdminModels = make(map[string]interface{})
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

func (am *GowaManager) AddModel(table_name string, columns []string, model interface{}){
	var gowaTable GowaTable

	gowaTable.Title = table_name
	gowaTable.Columns = columns

	am.AdminTables[table_name] = gowaTable
	am.AdminModels[table_name] = model
}

func (am *GowaManager) RemoveModel(table_name string){
	delete(am.AdminTables, table_name)
	delete(am.AdminModels, table_name)
}

func (am *GowaManager) GetRoutes() goServerUtils.Routes {
	routes := goServerUtils.Routes{
		goServerUtils.Route{
			"GetTables",
			"GET",
			"/gowa/api/rest/tables",
			GetTables,
		},
		goServerUtils.Route{
			"ShowTable",
			"GET",
			"/gowa/api/rest/tables/show/{table}",
			ShowTable,
		},
	};
	return routes
}


