package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type People struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
}

var peopleRepository []People 

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "1234"
	dbName := "peoples"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetPeoples(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	selDB, err := db.Query("SELECT * FROM peoples")

	if err != nil {
		panic(err.Error())
	}

	var peoples []People

	for selDB.Next() {
		
		var id int
		var firstname, lastname string

		err = selDB.Scan(&id, &firstname, &lastname)
		if err != nil {
			panic(err.Error())
		}

		peoples = append(peoples, People{ID: string(id), Firstname: firstname, Lastname: lastname})

	}	

    json.NewEncoder(w).Encode(peoples)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range peopleRepository {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&People{}) 
}

func CreatePeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var person People 
    _ = json.NewDecoder(r.Body).Decode(&person) 
    person.ID = params["id"] 
    peopleRepository = append(peopleRepository, person) 
    json.NewEncoder(w).Encode(peopleRepository)
}

func DeletePeople(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range peopleRepository {
        if item.ID == params["id"] {
            peopleRepository = append(peopleRepository[:index], peopleRepository[index+1:]...)
            break
        }
        json.NewEncoder(w).Encode(peopleRepository)
    }
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/contato", GetPeoples).Methods("GET")
    router.HandleFunc("/contato/{id}", GetPeople).Methods("GET")
    router.HandleFunc("/contato/{id}", CreatePeople).Methods("POST")
    router.HandleFunc("/contato/{id}", DeletePeople).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}