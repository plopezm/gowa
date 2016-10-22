package gowa

import (
	"runtime"
	"path"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/plopezm/goServerUtils"
)

func GowaStart(dbtype string, dbpath string, pageSize uint32) *GowaManager{
	GM = new(GowaManager)
	GM.Init(dbtype, dbpath, pageSize)
	return GM
}

func GowaAddRoutes(router *mux.Router) *mux.Router{
	for _,route := range GM.GetRoutes(){
		var handler http.Handler
		handler = route.HandlerFunc
		handler = goServerUtils.Logger(handler, route.Name)
		fmt.Println("Adding route: "+route.Name+" -> "+route.Method+" "+route.Pattern);
		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router;
}

func GowaGetTemplatePath() string{
	_, filename, _, _ := runtime.Caller(0)
	//fmt.Println(path.Join(path.Dir(filename), "template"))
	return path.Join(path.Dir(filename), "template")
}
