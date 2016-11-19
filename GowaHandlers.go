package gowa

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"reflect"
	"github.com/gorilla/sessions"
	"encoding/gob"
	"encoding/base64"
	"crypto/rand"
	"errors"
)

var store * sessions.CookieStore

func init(){
	bytes, _ := generateRandomBytes(64);
	//Registers user in gob to use as session variable
	store = sessions.NewCookieStore(bytes)
	gob.Register(&GowaUser{});
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func validateSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, *GowaUser, error){
	var user *GowaUser
	var ok bool

	session, err := store.Get(r, "gowasession")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, nil, err
	}

	val := session.Values["user"];
	if user, ok = val.(*GowaUser); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, nil, errors.New("User not found")
	}

	return session, user, nil

}

/****************************************
 *     LOGIN MANAGEMENT HANDLERS	*
 ****************************************/

func Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	usr, pass, ok := r.BasicAuth();
	if !ok {
		w.WriteHeader(http.StatusUnauthorized);
		return;
	}

	var user GowaUser;

	db, _:= GM.getSession();

	if db.Where("Email = ? AND Passwd = ?", usr, pass).First(&user).RecordNotFound(){
		w.WriteHeader(http.StatusNotFound);
		return;
	}

	session, _ := store.Get(r, "gowasession")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	session.Values["user"] = user
	session.Save(r, w);
	w.WriteHeader(http.StatusOK);
}

func ValidateSession(w http.ResponseWriter, r *http.Request){
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}
	w.WriteHeader(http.StatusOK);
}

func CreateUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	session, err := store.Get(r, "gowasession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var user *GowaUser
	var ok bool

	val := session.Values["user"];
	if user, ok = val.(*GowaUser); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user != nil && user.Permission == PERM_RW {
		http.Error(w, "User not valid", http.StatusUnauthorized)
		return
	}

	var newUser GowaUser

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = r.Body.Close(); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &newUser); err != nil {
		http.Error(w, "Error: Unmarshalling Driver user", http.StatusBadRequest)
		return;
	}

	if !newUser.IsValid() {
		http.Error(w, "Error: User not valid", http.StatusBadRequest)
		return
	}

	newUser.Create()

	w.WriteHeader(http.StatusCreated)
}


/****************************************
 *     TABLE MANAGEMENT HANDLERS	*
 ****************************************/

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
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	var gowaTables []GowaTable
	var gowatable GowaTable

	for _, gowatable = range GM.adminTables {
		gowaTables = append(gowaTables, gowatable)
	}
	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(gowaTables); err != nil {
		panic(err)
	}
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	vars := mux.Vars(r);
	table := vars["table"];

	gowatable := GM.adminTables[table]

	db, _:= GM.getSession()

	typ := GM.adminTables[table].Model
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sliceType := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	slice := reflect.New(sliceType.Type()).Interface()

	db.Table(gowatable.Title).Limit(GM.PageSize).Find(slice)
	gowatable.Rows = slice

	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(gowatable); err != nil {
		panic(err)
	}
}

func ShowTableRows(w http.ResponseWriter, r *http.Request) {
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");

	vars := mux.Vars(r);
	table := vars["table"];

	db, _:= GM.getSession()

	typ := GM.adminTables[table].Model
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sliceType := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	slice := reflect.New(sliceType.Type()).Interface()

	db.Table(table).Limit(GM.PageSize).Find(slice)

	w.WriteHeader(http.StatusOK);
	if err := json.NewEncoder(w).Encode(slice); err != nil {
		panic(err)
	}
}

func AddTableRow(w http.ResponseWriter, r *http.Request) {
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	vars := mux.Vars(r)
	table := vars["table"]

	obj, err := deserialize(r, GM.adminTables[table].Model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db, _:= GM.getSession()
	if db.Table(table).Create(obj).RowsAffected == 0{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated);
}

func RemoveTableRow(w http.ResponseWriter, r *http.Request) {
	_, _, err := validateSession(w, r)
	if(err != nil) {
		w.Write([]byte(err.Error()));
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	vars := mux.Vars(r)
	table := vars["table"]

	obj, err := deserialize(r, GM.adminTables[table].Model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db, _:= GM.getSession()
	if db.Table(table).Delete(obj).RowsAffected == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK);
}

