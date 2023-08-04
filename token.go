package main

import (
	"log"
	"strings"
)

func textToPostingLists(invertedIndex InvertedIndex, doc *Document) error {
	for i, t := range tokenize(doc.Body) {
		err := tokenToPostingList(invertedIndex, t, doc.ID, int64(i))
		if err != nil {
			log.Println("Fail to token to PostingList")
			return err
		}
	}
	return nil
}

func tokenize(s string) []string {
	// TODO: 色々拡張できる
	return strings.Split(strings.ToLower(s), " ")
}

func tokenToPostingList(invertedIndex InvertedIndex, token string, docID int64, pos int64) error {
	// 転置インデックスにtokenがあれば追加、なければpostingを新規作成
	if postings, ok := invertedIndex[token]; ok {
		isNewDoc := true
		for _, p := range postings {
			if docID == p.DocID {
				p.Positions = append(p.Positions, pos)
				isNewDoc = false
			}
		}
		if isNewDoc {
			postings = append(postings, &Posting{DocID: docID, Positions: []int64{pos}})
		}
		invertedIndex[token] = postings
	} else {
		postings := []*Posting{
			{DocID: docID, Positions: []int64{pos}},
		}
		invertedIndex[token] = postings
	}
	return nil
}
