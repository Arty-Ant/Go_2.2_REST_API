package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	// Фактически, мы вернём значение из интервала [min, min + (max - min + 1)]
	return min + rand.Int63n(max-min+1)
}

type RandomAccountParams struct {
	Owner    string `faker:"first_name"`
	Balance  int64
	Currency string `faker:"oneof: USD, EUR"`
}

func RandomAccount() RandomAccountParams {
	rap := RandomAccountParams{}
	err := faker.FakeData(&rap)
	if err != nil {
		log.Fatal(err)
	}
	return rap
}

type RandomUserParams struct {
	Username       string `faker:"username"`
	HashedPassword string `faker:"password"`
	FullName       string `faker:"name"`
	Email          string `faker:"email"`
}

func RandomUser() RandomUserParams {
	rup := RandomUserParams{}
	err := faker.FakeData(&rup)
	if err != nil {
		log.Fatal(err)
	}
	return rup
}
