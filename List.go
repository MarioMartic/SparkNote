package main

import (
	"database/sql"
	"fmt"
	"log"
)

type List struct {
	ID 			int 	`json:"id"`
	Title		string  `json:"title"`
	UserID		int		`json:"user_id"`
	Items	    []Item	`json:"items"`
	CreatedAt	string	`json:"created_at"`
}

type Item struct {
	ID 			int 	`json:"id"`
	Body		string  `json:"body"`
	Done		bool	`json:"done"`
	ListID		int		`json:"-"`
	CreatedAt	string	`json:"created_at"`
}

func (l *List) createList(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO list(title, user_id, deleted) VALUES(?, ?, ?)")
	_, err := db.Exec(statement, l.Title, l.UserID, 0)

	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&l.ID)

	if err != nil {
		return err
	}

	return nil
}

func (l *List) getList(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT title, user_id, created_at FROM list WHERE id=? AND deleted = 0")
	return db.QueryRow(statement, l.ID).Scan(&l.Title, &l.UserID, &l.CreatedAt)
}

func (l *List) updateList(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE list SET title=? WHERE id=?")
	_, err := db.Exec(statement, l.Title, l.ID)
	return err
}

func (l *List) deleteList(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE list SET deleted=1 WHERE id=?")
	_, err := db.Exec(statement, l.ID)
	return err
}

func getLists(db *sql.DB, userID int) ([]List, error) {

	statement := fmt.Sprintf("SELECT id, title, created_at FROM list WHERE user_id =? AND deleted = 0")
	rows, err := db.Query(statement, userID)

	defer rows.Close()

	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	lists := []List{}

	for rows.Next() {

		var l List

		err := rows.Scan(&l.ID, &l.Title, &l.CreatedAt)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		err = l.getListItems(db)

		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		lists = append(lists, l)
	}

	if err != nil {
		fmt.Sprintf(err.Error())
		return nil, err
	}

	return lists, nil
}

func (l *List) getListItems(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT title, user_id, created_at FROM list WHERE id =? AND deleted = 0")
	err := db.QueryRow(statement, l.ID).Scan(&l.Title, &l.UserID, &l.CreatedAt)

	if err != nil {
		log.Printf(err.Error())
		return err
	}

	statement = fmt.Sprintf("SELECT id, body, done, list_id, created_at FROM item WHERE list_id = ?")
	items, err := db.Query(statement, l.ID)

	defer items.Close()

	if err != nil {
		fmt.Sprintf(err.Error())
		return err
	}

	for counter := 0; items.Next(); counter++ {

		var listItem Item

		err := items.Scan(&listItem.ID, &listItem.Body, &listItem.Done, &listItem.ListID, &listItem.CreatedAt)

		if err != nil {
			log.Fatal(err)
			return err
		}

		l.Items = append(l.Items, listItem)

	}

	return nil
}


func (i *Item) createItemForList(db *sql.DB, listID int) error {
	statement := fmt.Sprintf("INSERT INTO item(body, done, list_id) VALUES(?, ?, ?)")
	_, err := db.Exec(statement, i.Body, 0, listID)

	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *Item) getItem(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT body, done, list_id, created_at FROM item WHERE id=?")
	return db.QueryRow(statement, i.ID).Scan(&i.Body, &i.Done, &i.ListID, &i.CreatedAt)
}

func (i *Item) updateItem(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE item SET body=?, done=? WHERE id=?")
	_, err := db.Exec(statement, i.Body, i.Done, i.ID)
	return err
}

func (i *Item) deleteItem(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM item WHERE id=?")
	_, err := db.Exec(statement, i.ID)
	return err
}