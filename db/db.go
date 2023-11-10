package db

import (
	"fmt"
	"sync"

	"github.com/gocql/gocql"
)

var cassandraOnce sync.Once
var cassandraDB *gocql.Session

func InitDatabase() (err error) {
	cassandraOnce.Do(func() {
		cluster := gocql.NewCluster("localhost")
		keyspace := "ganjoor"
		cluster.Keyspace = keyspace
		cluster.Port = 9042

		var session *gocql.Session
		session, err = cluster.CreateSession()
		if err != nil {
			return
		}

		cassandraDB = session

		query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}", keyspace)
		if err := session.Query(query).Exec(); err != nil {
			fmt.Println("Error creating keyspace:", err)
			return
		}

		fmt.Printf("Keyspace '%s' created (if not exists) successfully\n", keyspace)
	})

	if err != nil {
		return
	}

	return nil
}

func GetDB() (*gocql.Session, error) {
	if cassandraDB == nil {
		err := InitDatabase()
		if err != nil {
			return nil, err
		}
	}
	return cassandraDB, nil
}
