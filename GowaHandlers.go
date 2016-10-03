package gowa

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
)


func GetTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	var gowaTables []GowaTable
	var gowatable GowaTable

	db, _:= GM.GetSession()

	for _, gowatable = range GM.AdminTables {

		db.Table(gowatable.Title).Limit(GM.pageSize).Find(GM.AdminModels[gowatable.Title])
		gowatable.Rows = GM.AdminModels[gowatable.Title]

		gowaTables = append(gowaTables, gowatable)

	}
	fmt.Println(gowaTables)
	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(gowaTables); err != nil {
		panic(err)
	}
}


func ShowTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	vars := mux.Vars(r);
	table := vars["table"];

	db, _:= GM.GetSession()

	db.Table(table).Limit(GM.pageSize).Find(GM.AdminModels[table])

	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(GM.AdminModels[table]); err != nil {
		panic(err)
	}
}
