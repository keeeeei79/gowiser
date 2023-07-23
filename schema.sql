DROP TABLE document;
CREATE TABLE document (
	"id" SERIAL PRIMARY KEY,
	"title" text NOT NULL,
	"body" text NOT NULL
);


DROP TABLE inverted_index;
CREATE TABLE inverted_index (
	"id" SERIAL PRIMARY KEY,
	"token" text NOT NULL,
    "doc_id" integer NOT NULL,
    "positions" integer[] NOT NULL,
    UNIQUE(token, doc_id)
);