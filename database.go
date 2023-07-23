package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DBInvertedIndex struct {
	ID        int64   `db:"id"`
	Token     string  `db:"token"`
	DocID     int64   `db:"doc_id"`
	Positions []int64 `db:"positions"`
}

func dbAddDocument(doc *Document) (int64, error) {
	db, err := sqlx.Connect("pgx", "host=127.0.0.1 port=35432 dbname=gowiser user=gowiser password=gowiser sslmode=disable")
	defer db.Close()
	if err != nil {
		return 0, err
	}

	var id int64
	query := `INSERT INTO document (title, body) VALUES ($1, $2) RETURNING id`
	err = db.QueryRowx(query, doc.Title, doc.Body).Scan(&id)
	if err != nil {
		return 0, err
	}
	fmt.Println("User inserted successfully with ID: ", id)
	return id, nil

}

func dbUpsertPosting(token string, posting *Posting) error {
	db, err := sqlx.Connect("pgx", "host=127.0.0.1 port=35432 dbname=gowiser user=gowiser password=gowiser sslmode=disable")
	defer db.Close()
	if err != nil {
		return err
	}
	query := `SELECT * FROM public.inverted_index WHERE token = $1 and doc_id = $2`
	var dbInvertedIndex DBInvertedIndex
	err = db.Get(&dbInvertedIndex, query, token, posting.DocID)
	if err == nil {
		// もし存在していたらpositionsに追加して更新する
		dbInvertedIndex.Positions = append(dbInvertedIndex.Positions, posting.Positions...)
		_, err = db.Exec(`UPDATE inverted_index SET postings=$1 WHERE id = $2 `, dbInvertedIndex.Positions, dbInvertedIndex.ID)
		if err != nil {
			fmt.Println("Fail to update inverted index")
			return err
		}
		return nil

	}
	if err != sql.ErrNoRows {
		// A real error occurred
		return err
	}
	// 存在していない場合はInsert
	dbInvertedIndex = DBInvertedIndex{
		Token:     token,
		DocID:     posting.DocID,
		Positions: posting.Positions,
	}
	insertQuery := `INSERT INTO inverted_index (token, doc_id, positions) VALUES ($1, $2, $3)`
	_, err = db.Exec(insertQuery, token, posting.DocID, posting.Positions)
	if err != nil {
		fmt.Println("Fail to insert inverted index")
		return err
	}
	return nil
}
