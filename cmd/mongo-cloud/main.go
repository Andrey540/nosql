package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strings"
	"sync"
	"time"
)

const threadsCount = 10
const rowsPerThread = 1000

func main() {
	const DB_NAME = "db1"
	DB_HOSTS := []string{
		"rc1a-3ft122hmu8yngl2k.mdb.yandexcloud.net:27018",
	}
	const DB_USER = "user1"
	const DB_PASS = "lr8akUjs>"

	const CACERT = "/home/andreyegoshin/.mongodb/root.crt"

	url := fmt.Sprintf("mongodb://%s:%s@%s/%s?tls=true&tlsCaFile=%s",
		DB_USER,
		DB_PASS,
		strings.Join(DB_HOSTS, ","),
		DB_NAME,
		CACERT)

	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return
	}

	collection := client.Database("db1").Collection("users")

	testWrites(*collection)
	testReads(*collection)
}

func testWrites(collection mongo.Collection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)

				ctx := context.Background()
				_, err := collection.InsertOne(ctx, bson.D{{"login", login}, {"name", "User"}})

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

func testReads(collection mongo.Collection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		go func(index int) {
			defer wg.Done()

			var result struct {
				name string
			}

			for i := 1; i <= rowsPerThread; i++ {
				login := fmt.Sprintf("login-%d-%d", index, i)

				filter := bson.D{{"login", login}}
				ctx := context.Background()
				err := collection.FindOne(ctx, filter).Decode(&result)
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
