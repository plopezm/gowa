package gowa

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"reflect"
)

func deserialize(r *http.Request, typ reflect.Type) (interface{}, error){

	data,e:=ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if e !=nil{
		return nil, e
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	payload:=reflect.New(typ).Interface()

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	return payload,nil

}

func GetTablesStruct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	var gowaTables []GowaTable
	var gowatable GowaTable

	//db, _:= GM.GetSession()

	for _, gowatable = range GM.AdminTables {

		//db.Table(gowatable.Title).Limit(GM.pageSize).Find(GM.AdminModels[gowatable.Title])
		//gowatable.Rows = GM.AdminModels[gowatable.Title]
		gowaTables = append(gowaTables, gowatable)

	}
	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(gowaTables); err != nil {
		panic(err)
	}
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	vars := mux.Vars(r);
	table := vars["table"];

	gowatable := GM.AdminTables[table]

	db, _:= GM.GetSession()

	typ := GM.AdminTables[table].Model
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sliceType := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	slice := reflect.New(sliceType.Type()).Interface()

	db.Table(gowatable.Title).Limit(GM.pageSize).Find(slice)
	gowatable.Rows = slice

	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(gowatable); err != nil {
		panic(err)
	}
}

func ShowTableRows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	vars := mux.Vars(r);
	table := vars["table"];

	db, _:= GM.GetSession()

	typ := GM.AdminTables[table].Model
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sliceType := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	slice := reflect.New(sliceType.Type()).Interface()

	db.Table(table).Limit(GM.pageSize).Find(slice)

	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(slice); err != nil {
		panic(err)
	}
}

func AddTableRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	vars := mux.Vars(r)
	table := vars["table"]

	obj, err := deserialize(r, GM.AdminTables[table].Model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db, _:= GM.GetSession()
	if db.Table(table).Create(obj).RowsAffected == 0{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated);
}

func RemoveTableRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	vars := mux.Vars(r)
	table := vars["table"]

	obj, err := deserialize(r, GM.AdminTables[table].Model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db, _:= GM.GetSession()
	if db.Table(table).Delete(obj).RowsAffected == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK);
}

