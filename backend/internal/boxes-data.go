package boxes

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Box struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Likes   int64  `json:"likes"`
}

type BoxRepo struct {
	DB *sql.DB
}

// Get box by ID
func (r *BoxRepo) BoxRetrieve(id int) (*Box, error) {
	fail := func(err error) (*Box, error) {
		return nil, fmt.Errorf("BoxRetrieve: %v", err)
	}

	box := &Box{}
	row := r.DB.QueryRow("SELECT * FROM box WHERE id = ?", id)
	err := row.Scan(&box.ID, &box.Content, &box.Author, &box.Likes)
	if err != nil {
		return fail(err)
	}
	return box, nil
}

// Query all boxes
func (r *BoxRepo) BoxesList() ([]Box, error) {
	fail := func(err error) ([]Box, error) {
		return nil, fmt.Errorf("BoxesList: %v", err)
	}

	var boxes []Box
	rows, err := r.DB.Query("SELECT * FROM box")
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	for rows.Next() {
		var box Box

		if err := rows.Scan(
			&box.ID, &box.Content, &box.Author, &box.Likes,
		); err != nil {
			return fail(err)
		}
		boxes = append(boxes, box)
	}
	return boxes, err
}

// Insert new box, also updating outbox
func (r *BoxRepo) BoxInsert(box *Box) error {
	fail := func(err error) error {
		return fmt.Errorf("BoxInsert: %v", err)
	}

	// Start transaction with insert to box and outbox tables
	tx, err := r.DB.Begin()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	res, err := tx.Exec("INSERT INTO box (content, author, likes) VALUES (?, ?, ?)", box.Content, box.Author, box.Likes)
	if err != nil {
		return fail(err)
	}

	// Update ID
	id, err := res.LastInsertId()
	if err != nil {
		return fail(err)
	}
	box.ID = id

	// Make outbox
	json, err := json.Marshal(NewBoxEvent{BoxID: box.ID, Content: box.Content, Author: box.Author})
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO outbox (processed, type, payload) VALUES (FALSE, 'boxNew', ?)",
		json,
	)
	if err != nil {
		return fail(err)
	}

	// Fin
	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return nil
}

// Update liked box
func (r *BoxRepo) BoxLikeUpdate(id int64) error {
	fail := func(err error) error {
		return fmt.Errorf("BoxLikeUpdate: %v", err)
	}

	// Start transaction with insert to box and outbox tables
	tx, err := r.DB.Begin()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	_, err = tx.Exec("UPDATE box SET likes = likes + 1 WHERE id = ?", id)
	if err != nil {
		return fail(err)
	}

	row := tx.QueryRow("SELECT author, likes FROM box WHERE id = ?", id)
	var likes int64
	var author string
	err = row.Scan(&author, &likes)
	if err != nil {
		return fail(err)
	}

	// Make outbox
	json, err := json.Marshal(BoxLikedEvent{BoxID: id, Author: author, Likes: likes})
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO outbox (processed, type, payload) VALUES (FALSE, 'boxLike', ?)",
		json,
	)
	if err != nil {
		return fail(err)
	}

	// Fin
	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return nil
}
