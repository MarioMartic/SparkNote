package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"log"
)

func (a *App) getNotes(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	notes, err := getNotes(a.DB, user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, notes)
}

func (a *App) createNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		log.Println(err.Error())
		return
	}

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		log.Println(err.Error())
		return
	}

	n.UserID = user.ID

	defer r.Body.Close()

	if err := n.createNote(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		log.Println(err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, n)
}

func (a *App) getNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Note ID")

		return
	}

	note := Note{ID: id}

	if err := note.getNote(a.DB); err != nil {
		switch err {

		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Note not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	user, err := getUserFromToken(w, r)

	if note.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access!")
		return
	}

	respondWithJSON(w, http.StatusOK, note)

}

func (a *App) updateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Note ID")
		return
	}

	var n Note
	n.ID = id

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&n); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := n.updateNote(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, n)
}

func (a *App) deleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	log.Print(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Note ID")
		return
	}

	user, err := getUserFromToken(w, r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid User!")
		return
	}

	n := Note{ID: id}

	n.getNote(a.DB)

	if n.UserID != user.ID {
		respondWithError(w, http.StatusForbidden, "Unauthorized access!")
		log.Print(n.UserID, user.Username)
		return
	}

	if err := n.deleteNote(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		log.Print(err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}