package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"sync"

	"time"
)

const threadsCount = 10
const rowsPerThread = 100000

func main() {
	cluster := gocql.NewCluster("172.32.0.2", "172.32.0.3", "172.32.0.4")
	cluster.Keyspace = "system"
	cluster.Consistency = gocql.Any
	// connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS test_keyspace WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	err = session.Query("CREATE TABLE IF NOT EXISTS test_keyspace.user (id uuid, login text, created_at timestamp, PRIMARY KEY (id));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	testWrites(*session)
	testReads(*session)
}

func testWrites(session gocql.Session) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)
				uuid, err := gocql.RandomUUID()
				if err != nil {
					log.Println(err)
					return
				}
				err = session.Query("INSERT INTO test_keyspace.user (id, login, created_at) VALUES (?, ?, '2018-01-07');", uuid, login).Exec()
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(x)
	}
	wg.Wait()
	fmt.Println("Write time:", time.Now().UnixMilli()-currentTime)
}

func testReads(session gocql.Session) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)
				scanner := session.Query("SELECT id, login, created_at FROM user WHERE login = ?", login).Iter().Scanner()
				for scanner.Next() {
				}
			}
		}(x)
	}
	wg.Wait()
	fmt.Println("Read time:", time.Now().UnixMilli()-currentTime)
}
