package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"runtime"

	"time"
)

func main() {
	cluster := gocql.NewCluster("172.22.0.2", "172.22.0.3", "172.22.0.4")
	cluster.Keyspace = "test_keyspace"
	cluster.Consistency = gocql.Quorum
	// connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	currentTime := time.Now().UnixMilli()
	var currentMemory runtime.MemStats
	runtime.ReadMemStats(&currentMemory)

	scanner := session.Query(`SELECT id, login, created_at FROM user`).Iter().Scanner()
	for scanner.Next() {
		var (
			id        gocql.UUID
			login     string
			createdAt time.Time
		)
		err = scanner.Scan(&id, &login, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("User:", id, login, createdAt)
	}
	// scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	fmt.Println("Time:", time.Now().UnixMilli()-currentTime)
	fmt.Println("Memory:", (memory.Alloc-currentMemory.Alloc)/1024)
}
