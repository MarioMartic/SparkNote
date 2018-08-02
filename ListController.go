package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
)

func (a *App) getLists(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := getLists(a.DB, user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, lists)
}

func (a *App) createList(w http.ResponseWriter, r *http.Request) {
	var l List
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	l.UserID = user.ID

	defer r.Body.Close()

	if err := l.createList(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, l)
}

func (a *App) updateList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid List ID")
		return
	}

	l := List{ID: id}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := l.updateList(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, l)
}

func (a *App) deleteList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var l List
	l.ID = id

	if err = l.getList(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}


	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if l.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access!")
		return
	}


	if err := l.deleteList(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getListItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid List ID")
		return
	}

	list := List{ID: id}

	if err := list.getListItems(a.DB); err != nil {
		switch err {

		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "List items not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	user, err := getUserFromToken(w, r)

	if list.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access!")
		return
	}

	respondWithJSON(w, http.StatusOK, list.Items)

}

func (a *App) createListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	l := List{ID: listID}

	if err = l.getList(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't get list!")
		return
	}

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	if l.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access!")
		return
	}

	i := Item{ListID: listID}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err = i.createItemForList(a.DB, i.ListID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't get item")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateOrDeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	itemID, err := strconv.Atoi(vars["item"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var l List
	l.ID = id

	if err = l.getList(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't get list!")
		return
	}

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	i := Item{ID: itemID}

	if err = i.getItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't get item")
	}

	if l.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access!")
		return
	}

	switch r.Method {

	case "DELETE":
		if err := i.deleteItem(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, "DELETE ERROR")
			return
		}
	case "PUT":

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&i); err != nil {
			respondWithError(w, http.StatusBadRequest, "DECODER ERROR")
			return
		}

		defer r.Body.Close()

		if err := i.updateItem(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, "UPDATE ERROR")
			return
		}
	default:
		return

	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
