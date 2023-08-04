package main

import (
	"log"
	"sort"

	"github.com/jmoiron/sqlx"
)


func searchDocs(db *sqlx.DB, tokens []string) ([]*Document, error) {
	candidates, err := tokenSearch(db, tokens)
	if err != nil {
		log.Println("Fail to tokenSearch: ", err)
		return nil, err
	}
	res, err := intersect(candidates)
	if err != nil {
		log.Println("Fail to intersect: ", err)
		return nil, err
	}
	return res, nil
}

func tokenSearch(db *sqlx.DB, tokens []string) ([][]*Document, error) {
	candidates := make([][]*Document, 0)
	for _, t := range tokens {
		docs, err := dbSearch(db, t)
		if err != nil {
			log.Println("Fail to dbSearch: ", err)
			return nil, err
		}
		candidates = append(candidates, docs)
	}
	return candidates, nil

}

func intersect(candidates [][]*Document) ([]*Document, error){
	// 1番小さい候補集合の各ドキュメントを調べていく
	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i]) < len(candidates[j])
	})

	curs := make([]int64, len(candidates))
	isDone := false
	res := make([]*Document, 0)
	// wiserに従ってAとする
	A := candidates[0]
	for i:=0; i < len(A); i++ {
		docID := A[i].ID
		isExist := true
		// 他の候補集合にそのドキュメントが存在するかを調べる
		for j:= 1; j<len(candidates); j++ {
			ok, cur := checkExist(docID, candidates[j], curs[j])
			curs[j] = cur
			if (cur >= int64(len(candidates[j]))) {
				isDone = true
			}
			if (!ok) {
				isExist = ok
				break
			}
		}
		if (isExist) {
			res = append(res, A[i])
		}
		if (isDone) {
			break
		}
	}
	return res, nil
}

func checkExist(docID int64, docs []*Document, cur int64) (bool, int64){
	isExist := false
	for {
		if (cur >= int64(len(docs))) {
			break
		}
		if (docID == docs[cur].ID) {
			isExist = true
			cur++
			break
		}
		if (docID < docs[cur].ID ) {
			break
		}
		cur++
	}
	return isExist, cur
}