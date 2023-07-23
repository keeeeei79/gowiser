package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

const updateThreshold = 2

type InvertedIndex = map[string][]*Posting

func addDocs(db *sqlx.DB, docs []*Document) error {
	invertedIndex := make(map[string][]*Posting)
	for i, doc := range docs {
		docID, err := dbAddDocument(db, doc)
		if err != nil {
			log.Println("Fail to add a document: ", err)
			return err
		}
		doc.ID = docID
		// 文書をtokenに分割し、転置インデックスに格納する
		err = textToPostingLists(invertedIndex, doc)
		if err != nil {
			log.Println("Fail to convert to inverted index: ", err)
			return err
		}
		// 所定の文書数が貯まったらストレージ上の転置インデックスとマージして更新
		if i != 0 && i%updateThreshold == 0 {
			err = mergeInvertedIndex(db, invertedIndex)
			if err != nil {
				log.Println("Fail to merge with inverted index in db: ", err)
				return err
			}
			invertedIndex = make(map[string][]*Posting)
		}
	}
	if len(invertedIndex) > 0 {
		err := mergeInvertedIndex(db, invertedIndex)
		if err != nil {
			log.Println("Fail to merge with inverted index in db: ", err)
			return err
		}

	}
	return nil
}

func mergeInvertedIndex(db *sqlx.DB, invertedIndex InvertedIndex) error {
	for token, postings := range invertedIndex {
		for _, p := range postings {
			err := dbUpsertInvertedIndex(db, token, p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
