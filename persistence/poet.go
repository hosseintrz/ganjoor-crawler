package persistence

import (
	"github.com/gocql/gocql"
	"github.com/hosseintrz/ganjoor_crawler/model"
)

func InsertPoet(session *gocql.Session, poet *model.Poet) error {
	query := `
		INSERT INTO poets (id, name, description)
		VALUES (?, ?, ?)
	`

	return session.Query(
		query,
		poet.ID,
		poet.Name,
		poet.Description,
	).Exec()
}

func GetPoetByName(session *gocql.Session, name string) (*model.Poet, error) {
	var poet model.Poet
	query := "SELECT id, name, description FROM poets WHERE name = ? LIMIT 1"
	iter := session.Query(query, name).Iter()
	if iter.Scan(&poet.ID, &poet.Name, &poet.Description) {
		return &poet, nil
	}
	return nil, iter.Close()
}
