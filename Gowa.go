package gowa

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GowaStart(dbtype string, dbpath string, pageSize uint32) *GowaManager{
	GM = new(GowaManager)
	GM.Init(dbtype, dbpath, pageSize)
	return GM
}

func GowaEnableWebAdmin(router *mux.Router, webpath string){
	router.PathPrefix(webpath).Handler(http.FileServer(http.Dir("template")));
}
