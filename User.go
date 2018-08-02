package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID 			int 	`json:"id"`
	Username	string  `json:"username"`
	Password	[]byte	`json:"-"`
	Name 		string 	`json:"name"`
}

func (u *User) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name FROM users WHERE id=?")
	return db.QueryRow(statement, u.ID).Scan(&u.Name)
}

func (u *User) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET name=?, age=? WHERE id=?")
	_, err := db.Exec(statement, u.Name, u.ID)
	return err
}

func (u *User) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=?")
	_, err := db.Exec(statement, u.ID)
	return err
}

func (u *User) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO users(username, password, name) VALUES(?, ?, ?)")
	_, err := db.Exec(statement, u.Username, u.Password, u.Name)

	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) findByUsername(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, password, name FROM users WHERE username=?")
	return db.QueryRow(statement, u.Username).Scan(&u.ID, &u.Password, &u.Name)
}

func getUsers(db *sql.DB, start, count int) ([]User, error) {
	statement := fmt.Sprintf("SELECT id, name, username FROM users LIMIT ? OFFSET ?")
	rows, err := db.Query(statement, count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next(){
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}