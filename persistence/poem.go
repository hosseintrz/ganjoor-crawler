package persistence

import (
	"github.com/gocql/gocql"
	"github.com/hosseintrz/ganjoor_crawler/model"
)

func InsertPoem(session *gocql.Session, poem model.Poem) error {
	query := `
		INSERT INTO poems (id, poet_id, poet_name, title, content)
		VALUES (?, ?, ?, ?, ?)
	`

	return session.Query(
		query,
		poem.ID,
		poem.PoetID,
		poem.PoetName,
		poem.Title,
		poem.Content,
	).Exec()
}
