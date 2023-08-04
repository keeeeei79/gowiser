package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DBInvertedIndex struct {
	ID        int64   `db:"id"`
	Token     string  `db:"token"`
	DocID     int64   `db:"doc_id"`
	Positions []int64 `db:"positions"`
}

func dbAddDocument(db *sqlx.DB, doc *Document) (int64, error) {
	var id int64
	query := `INSERT INTO document (title, body) VALUES ($1, $2) RETURNING id`
	err := db.QueryRowx(query, doc.Title, doc.Body).Scan(&id)
	if err != nil {
		log.Println("Fail to insert document")
		return 0, err
	}
	fmt.Println("Document was inserted successfully with ID: ", id)
	return id, nil

}

func dbUpsertInvertedIndex(db *sqlx.DB, token string, posting *Posting) error {
	query := `SELECT * FROM public.inverted_index WHERE token = $1 and doc_id = $2`
	var dbInvertedIndex DBInvertedIndex
	err := db.Get(&dbInvertedIndex, query, token, posting.DocID)
	if err == nil {
		// もし存在していたらpositionsに追加して更新する
		dbInvertedIndex.Positions = append(dbInvertedIndex.Positions, posting.Positions...)
		_, err = db.Exec(`UPDATE inverted_index SET postings=$1 WHERE id = $2 `, dbInvertedIndex.Positions, dbInvertedIndex.ID)
		if err != nil {
			log.Println("Fail to update inverted index")
			return err
		}
		fmt.Println("Inverted Index was updated successfully with ID: ", dbInvertedIndex.ID)
		return nil

	}
	if err != sql.ErrNoRows {
		log.Println("Fail to db.Get")
		return err
	}
	// 存在していない場合はInsert
	insertQuery := `INSERT INTO inverted_index (token, doc_id, positions) VALUES ($1, $2, $3)`
	_, err = db.Exec(insertQuery, token, posting.DocID, posting.Positions)
	if err != nil {
		log.Println("Fail to insert inverted index")
		return err
	}
	fmt.Println(token, "was inserted to inverted index successfully")
	return nil
}

func dbSearch(db *sqlx.DB, keyword string) ([]*Document, error) {
	var docIDs []int64
	err := db.Select(&docIDs, "SELECT doc_id FROM public.inverted_index WHERE token = $1", keyword)
	if err != nil {
		log.Println("Fail to select from inverted index")
		return nil, err
	}
	fmt.Println(keyword, "Hit docs ids: ", docIDs)
	if (len(docIDs)== 0){
		return []*Document{}, nil
	}

	sql := "SELECT * FROM public.document WHERE id IN (?) ORDER BY id"
	sql, params, err := sqlx.In(sql, docIDs)
	if err != nil {
		log.Println("Fail to use IN")
		return nil, err
	}
	sql = db.Rebind(sql)

	var docs []*Document
	err = db.Select(&docs, sql, params...)
	if err != nil {
		log.Println("Fail to select from document")
		return nil, err
	}
	return docs, nil
}
