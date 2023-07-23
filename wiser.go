package main

import (
	"log"
)

const updateThreshold = 2

type InvertedIndex = map[string][]*Posting

func addDocs(docs []*Document) error {
	invertedIndex := make(map[string][]*Posting)
	for i, doc := range docs {
		docID, err := dbAddDocument(doc)
		if err != nil {
			log.Fatalln("Fail to add a document: ", err)
		}
		doc.ID = docID
		// 文書をtokenに分割し、転置インデックスに格納する
		err = textToPostingLists(invertedIndex, doc)
		if err != nil {
			log.Fatalln("Fail to convert to inverted index: ", err)
		}
		// 所定の文書数が貯まったらストレージ上の転置インデックスとマージして更新
		if i != 0 && i%updateThreshold == 0 {
			err = mergeInvertedIndex(invertedIndex)
			if err != nil {
				log.Fatalln("Fail to merge with inverted index in db: ", err)
			}
			invertedIndex = make(map[string][]*Posting)
		}
	}
	err := mergeInvertedIndex(invertedIndex)
	if err != nil {
		log.Fatalln("Fail to merge with inverted index in db: ", err)
	}
	return nil
}

func mergeInvertedIndex(invertedIndex InvertedIndex) error {
	for token, postings := range invertedIndex {
		for _, p := range postings {
			err := dbUpsertPosting(token, p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
