package main

import (
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"fmt"
	"database/sql"
	"log"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var a App

func main() {

	initKeys()

	a.initializeDB("root", "", "spark_note")
	a.initializeRoutes()

	a.Run(":8080")

}

func (a *App) initializeDB(user, password, dbname string) {

	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error

	a.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Now listening...")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, handler))
}

func optionsRequest(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
