package main

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/callicoder/go-docker/pkg/common/infrastructure/server"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	threadsCount  = 100
	rowsPerThread = 100000
	idKey         = "id"
	firstNameKey  = "first_name"
	lastNameKey   = "last_name"
	genderKey     = "gender"
	emailKey      = "email"
	phoneKey      = "phone"
	registeredKey = "registered"
	countryKey    = "country"
	avatarKey     = "avatar"
)

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
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})
	if client == nil {
		log.Fatal("redis client pointer = nil")
		return
	}
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	logger := server.InitLogger()
	errorLogger := server.InitErrorLogger()
	stopChan := make(chan struct{})
	server.ListenOSKillSignals(stopChan)
	serverHub := server.NewHub(stopChan)

	server.ServeHTTP(
		":8090",
		serverHub,
		getHandler(ctx, client),
		getHandlerName(ctx, client),
		logger, errorLogger)

	err = serverHub.Wait()
	if err != nil {
		errorLogger.Println(err)
	}

	//testWrites(ctx, client)
	//testReads(ctx, client)
}

func getHandler(ctx context.Context, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		id := rand.Int63n(threadsCount * rowsPerThread * 2)
		key := fmt.Sprintf("user:%d", id)
		_, err := client.Get(ctx, key).Result()
		if err != nil {
			//	fmt.Println(err)
			if !errors.Is(err, redis.Nil) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func getHandlerName(ctx context.Context, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		id := rand.Int63n(threadsCount * rowsPerThread * 2)
		key := fmt.Sprintf("user:%d", id)
		_, err := client.Get(ctx, key).Result()
		if err != nil {
			//	fmt.Println(err)
			if !errors.Is(err, redis.Nil) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}

func testWrites(ctx context.Context, client *redis.Client) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)
				key := fmt.Sprintf("%d", id)
				user := generateUser(id)
				userMap := map[string]interface{}{idKey: user.ID, firstNameKey: user.FirstName, lastNameKey: user.FirstName, genderKey: strconv.Itoa(user.Gender), emailKey: user.Email, phoneKey: user.Phone, registeredKey: user.Registered, countryKey: user.Country, avatarKey: user.Avatar}
				_, err := client.HSet(ctx, key, userMap).Result()
				if err != nil {
					log.Println(err, key)
					return
				}
			}
		}(x)
	}
	wg.Wait()
	fmt.Println("Write time:", time.Now().UnixMilli()-currentTime)
}

func testReads(ctx context.Context, client *redis.Client) {
	var wg sync.WaitGroup
	wg.Add(threadsCount)

	currentTime := time.Now().UnixMilli()
	for x := 1; x <= threadsCount; x++ {
		baseId := uint64((x - 1) * rowsPerThread)
		go func(index int) {
			defer wg.Done()

			for i := 1; i <= rowsPerThread; i++ {
				id := baseId + uint64(i)
				key := fmt.Sprintf("%d", id)

				_, err := client.HGetAll(ctx, key).Result()
				if err != nil {
					log.Println(err, key)
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
