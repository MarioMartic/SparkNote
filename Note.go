package main

import (
	"database/sql"
	"fmt"
	"log"
)

type Note struct {
	ID 			int 	`json:"id"`
	Title		string  `json:"title"`
	Body		string	`json:"body"`
	UserID 		int 	`json:"user_id"`
	Archived	int		`json:"archived"`
	CreatedAt	string	`json:"created_at"`
}


func (n *Note) getNote(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, title, body, user_id, archived, created_at FROM note WHERE id=? AND deleted = 0")
	return db.QueryRow(statement, n.ID).Scan(&n.ID, &n.Title, &n.Body, &n.UserID, &n.Archived, &n.CreatedAt)
}

func (n *Note) updateNote(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE note SET title=?, body=?, archived=? WHERE id=?")
	_, err := db.Exec(statement, n.Title, n.Body, n.Archived, n.ID)
	return err
}

func (n *Note) deleteNote(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE note SET deleted = 1 WHERE id=?")
	_, err := db.Exec(statement, n.ID)
	return err
}

func (n *Note) createNote(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO note(title, body, archived, user_id) VALUES(?, ?, ?, ?)")
	_, err := db.Exec(statement, n.Title, n.Body, n.Archived, n.UserID)

	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = db.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&n.ID)

	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func getNotes(db *sql.DB, userID int) ([]Note, error) {

	statement := fmt.Sprintf("SELECT id, title, body, archived, user_id, created_at FROM note WHERE user_id = ? AND deleted = 0 ORDER BY created_at DESC")
	rows, err := db.Query(statement, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := []Note{}

	for rows.Next(){
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.Archived, &n.UserID, &n.CreatedAt); err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	return notes, nil
}
