package model

import "github.com/gocql/gocql"

type Poem struct {
	ID       gocql.UUID `cql:"id"`
	Title    string     `cql:"title"`
	Content  string     `cql:"content"`
	PoetID   gocql.UUID `cql:"poet_id"`
	PoetName string     `cql:"poet_name"`
}

type Poet struct {
	ID          gocql.UUID `cql:"id"`
	Name        string     `cql:"name"`
	Description string     `cql:"description"`
}
