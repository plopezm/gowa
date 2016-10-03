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
}


func (am *GowaManager) Init(dbtype string, dbpath string) error{
	var err error;
	am.DB, err = gorm.Open(dbtype, dbpath);
	if err != nil {
		panic(err)
	}
	am.DB.SingularTable(true)
	defer am.DB.Close();

	//am.DB.SingularTable(true)

	if err != nil {
		return err;
	}
	am.DbPath = dbpath
	am.DbType = dbtype
	am.AdminTables = make(map[string]GowaTable)
	am.AdminModels = make(map[string]interface{})
	return nil
}

func (am *GowaManager) GetSession() (*gorm.DB, error){
	var err error;

	if am.DB.DB().Ping() != nil{
		am.DB, err = gorm.Open(am.DbType, am.DbPath);
		am.DB.SingularTable(true)
		return am.DB, err;
	}
	return am.DB, nil;
}

func (am *GowaManager) End(){
	am.DB.Close();
}

func (am *GowaManager) AddModel(table_name string, columns []string, model interface{}) error{
	//if objects == nil{
	//	return errors.New("objects cannot be null")
	//}

	var gowaTable GowaTable

	gowaTable.Title = table_name
	gowaTable.Columns = columns

	am.AdminTables[table_name] = gowaTable
	am.AdminModels[table_name] = model

	return nil
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


