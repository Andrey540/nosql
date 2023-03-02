package main

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"log"
	"sync"
	"time"
)

const threadsCount = 10
const rowsPerThread = 100000

type User struct {
	Name string `json:"name"`
}

func main() {
	username := "test"
	password := "administrator"

	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}

	cluster, err := gocb.Connect("couchbase://172.18.0.3", options)
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket("travel-sample")

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := bucket.Scope("tenant_agent_00").Collection("users")

	testWrites(*collection)
	testReads(*collection)
}

func testWrites(collection gocb.Collection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)

				_, err := collection.Upsert(login,
					User{
						Name: fmt.Sprintf("Name-%d-%d", index, i),
					}, nil)

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

func testReads(collection gocb.Collection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)

				_, err := collection.Get(login, nil)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(x)
	}
	wg.Wait()
	fmt.Println("Read time:", time.Now().UnixMilli()-currentTime)
}
