package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/callicoder/go-docker/pkg/common/infrastructure/server"
	"github.com/couchbase/gocb/v2"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const threadsCount = 100
const rowsPerThread = 100000

type User struct {
	ID         uint64 `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     int    `json:"gender"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Registered string `json:"registered"`
	Country    string `json:"country"`
	Avatar     string `json:"avatar"`
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

	cluster, err := gocb.Connect("couchbase://172.18.0.5", options)
	if err != nil {
		log.Fatal(err)
	}
	defer cluster.Close(nil)

	bucket := cluster.Bucket("travel-sample")

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := bucket.Scope("tenant_agent_01").Collection("users")

	logger := server.InitLogger()
	errorLogger := server.InitErrorLogger()
	stopChan := make(chan struct{})
	server.ListenOSKillSignals(stopChan)
	serverHub := server.NewHub(stopChan)

	server.ServeHTTP(
		":8090",
		serverHub,
		getHandler(*collection),
		getNameHandler(*cluster),
		logger, errorLogger)

	err = serverHub.Wait()
	if err != nil {
		errorLogger.Println(err)
	}

	//testWrites(*collection)
	//testReads(*collection)
}

func getHandler(collection gocb.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		id := rand.Int63n(threadsCount * rowsPerThread * 2)
		_, err := collection.Get(strconv.FormatUint(uint64(id), 10), nil)
		if err != nil {
			// fmt.Println(err)
			if !errors.Is(err, gocb.ErrDocumentNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func getNameHandler(cluster gocb.Cluster) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		gender := rand.Intn(1)
		name := randomdata.FirstName(gender)
		queryResult, err := cluster.Query(
			fmt.Sprintf("SELECT * FROM `travel-sample`.tenant_agent_01.users WHERE first_name = '"+name+"' limit 1000;"),
			&gocb.QueryOptions{},
		)
		if err != nil {
			// fmt.Println(err)
			if !errors.Is(err, gocb.ErrDocumentNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		for queryResult.Next() {
			var result interface{}
			err = queryResult.Row(&result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func testWrites(collection gocb.Collection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)
				_, err := collection.Upsert(strconv.FormatUint(id, 10), generateUser(id), nil)

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
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)
				_, err := collection.Get(strconv.FormatUint(id, 10), nil)
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

func generateUser(id uint64) User {
	gender := rand.Intn(1)
	pic := rand.Intn(35)
	return User{
		ID:         id,
		FirstName:  randomdata.FirstName(gender),
		LastName:   randomdata.LastName(),
		Gender:     gender,
		Email:      randomdata.Email(),
		Phone:      randomdata.PhoneNumber(),
		Registered: randomdata.FullDate(),
		Country:    randomdata.Country(randomdata.FullCountry),
		Avatar:     fmt.Sprintf("https://randomuser.me/api/portraits/%d.jpg", pic),
	}
}
