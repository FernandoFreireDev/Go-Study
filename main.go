package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)

type People struct {
    ID        string   `gorm:primary_key"`
    Firstname string
    Lastname  string
}

type Retorno struct {
    Code int
    Message string
}

func DbConn() (*gorm.DB, error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "1234"
	dbName := "peoples"

	db, err := gorm.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	return db, err
}

func GetPeoples(w http.ResponseWriter, r *http.Request) {

    db, err := DbConn()
    checkError(err)

    var peoples []People

    db.Table("peoples").Find(&peoples)
	
    json.NewEncoder(w).Encode(peoples)
    fmt.Print(peoples)

    defer db.Close()
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    db, err := DbConn()
    checkError(err)

    var people People
    var count int
    db.Where("id = ?", params["id"]).First(&people).Count(&count)

    if count == 0 {
        retorno, _ := json.Marshal(Retorno{http.StatusNotFound, "Usuário não encontrado"})
        http.Error(w, string(retorno), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(people)
}

func CreatePeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    db, err := DbConn()
    checkError(err)

    var people People
    _ = json.NewDecoder(r.Body).Decode(&people) 
    people.ID = params["id"] 

    db.Create(&people)

    json.NewEncoder(w).Encode(people)
}

func DeletePeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    db, err := DbConn()
    checkError(err)

    db.Where("id = ?", params["id"]).Delete(People{})

    json.NewEncoder(w).Encode(Retorno{http.StatusOK, "Usuário deletado"});
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/contato", GetPeoples).Methods("GET")
    router.HandleFunc("/contato/{id}", GetPeople).Methods("GET")
    router.HandleFunc("/contato/{id}", CreatePeople).Methods("POST")
    router.HandleFunc("/contato/{id}", DeletePeople).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}

func checkError(err error) {
    if err != nil {
		panic(err.Error())
	}
}