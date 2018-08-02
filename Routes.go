package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/rs/cors"
)

var handler http.Handler

func (a *App) initializeRoutes() {

	a.Router = mux.NewRouter()

	//PUBLIC ENDPOINTS
	a.Router.HandleFunc("/", homePage)
	a.Router.HandleFunc("/login", LoginHandler).Methods("POST")
	//a.Router.HandleFunc("/login", optionsRequest).Methods("OPTIONS")
	a.Router.HandleFunc("/signup", RegisterHandler).Methods("POST")

	//a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

	a.Router.Handle("/users", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.getUsers)),
	)).Methods("GET")

	/////////////////////
	// NOTE CRUD Routes
	/////////////////////

	a.Router.Handle("/note", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.getNotes)),
	)).Methods("GET")

	a.Router.Handle("/note", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.createNote)),
	)).Methods("POST")

	a.Router.Handle("/note/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.getNote)),
	)).Methods("GET")

	a.Router.Handle("/note/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.updateNote)),
	)).Methods("PUT")

	a.Router.Handle("/note/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.deleteNote)),
	)).Methods("DELETE")


	/////////////////////
	// LIST CRUD Routes
	/////////////////////

	a.Router.Handle("/list", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.getLists)),
	)).Methods("GET")

	a.Router.Handle("/list", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.createList)),
	)).Methods("POST")

	a.Router.Handle("/list/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.updateList)),
	)).Methods("PUT")

	a.Router.Handle("/list/{id:[0-9]+}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.deleteList)),
	)).Methods("DELETE")


	///////////////////////////
	// LIST ITEMS CRUD Routes
	///////////////////////////

	a.Router.Handle("/list/{id:[0-9]+}/item", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.createListItem)),
	)).Methods("POST")

	a.Router.Handle("/list/{id:[0-9]+}/item", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.getListItems)),
	)).Methods("GET")


	a.Router.PathPrefix("/list/{id:[0-9]+}").Subrouter().Handle("/item/{item:[0-9]+}",
		negroni.New(
			negroni.HandlerFunc(ValidateTokenMiddleware),
			negroni.Wrap(http.HandlerFunc(a.updateOrDeleteItem)),
		)).Methods("PUT", "DELETE")

	handler = cors.AllowAll().Handler(a.Router)

}
