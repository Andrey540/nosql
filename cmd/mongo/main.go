package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"
)

const threadsCount = 10
const rowsPerThread = 100000

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.10.1:27017"))

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return
	}

	collection := client.Database("testing").Collection("users")

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
