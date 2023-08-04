package main

import (
	"reflect"
	"testing"
)

func TestIntersect(t *testing.T) {
	doc1 := &Document{ID: 1, Title: "A", Body: "Body A"}
	doc2 := &Document{ID: 2, Title: "B", Body: "Body B"}
	doc3 := &Document{ID: 3, Title: "C", Body: "Body C"}
	doc4 := &Document{ID: 4, Title: "D", Body: "Body D"}
	doc5 := &Document{ID: 5, Title: "E", Body: "Body E"}

	candidates := [][]*Document{
		{doc1, doc2, doc3, doc4, doc5},
		{doc3, doc5},
	}

	result, err := intersect(candidates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []*Document{doc3, doc5}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
}