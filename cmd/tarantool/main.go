package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/callicoder/go-docker/pkg/common/infrastructure/server"
	"github.com/tarantool/go-tarantool"
	"io"
	"log"
	"math/rand"
	"net/http"
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
	opts := tarantool.Opts{User: "guest", Pass: ""}
	address := fmt.Sprintf("%s:%s", "127.0.0.1", "3301")
	connection, err := tarantool.Connect(address, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	logger := server.InitLogger()
	errorLogger := server.InitErrorLogger()
	stopChan := make(chan struct{})
	server.ListenOSKillSignals(stopChan)
	serverHub := server.NewHub(stopChan)

	server.ServeHTTP(
		":8090",
		serverHub,
		getHandler(connection),
		getHandlerName(connection),
		logger, errorLogger)

	err = serverHub.Wait()
	if err != nil {
		errorLogger.Println(err)
	}

	//testWrites(connection)
	//testReads(connection)
}

func getHandler(connection *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		id := rand.Int63n(threadsCount * rowsPerThread * 2)
		var result []User
		err := connection.SelectTyped("user", "primary", 0, 1, tarantool.IterEq, []interface{}{id}, &result)
		if err != nil {
			//	fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func getHandlerName(connection *tarantool.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		gender := rand.Intn(1)
		name := randomdata.FirstName(gender)
		var result []User
		err := connection.SelectTyped("user", "first_name_idx", 0, 10000, tarantool.IterEq, []interface{}{name}, &result)
		if err != nil {
			//	fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func testWrites(connection *tarantool.Connection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)

				user := generateUser(id)
				var result []User
				err := connection.InsertTyped("user", []interface{}{user.ID, user.FirstName, user.LastName, user.Gender, user.Email, user.Phone, user.Registered, user.Phone, user.Avatar}, &result)

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

func testReads(connection *tarantool.Connection) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)
				var result []User
				err := connection.SelectTyped("user", "primary", 0, 1, tarantool.IterEq, []interface{}{id}, &result)
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
