package data

import (
	"database/sql"
	"fmt"

	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type Models struct {
	Upper db2.Session

	// any models inserted here (and in the New function)
	// are easily accessible throughout entire application
	Users         User
	Tokens        Token
	Subscriptions Subscription
	Products      Product
	Prices        Price
}

func New(dbPool *sql.DB) Models {
	upper, err := postgresql.New(dbPool)
	if err != nil {
		fmt.Printf("upper error -> %+v\n", err)
	}

	return Models{
		Upper: upper,
		//
		//
		Users:         User{},
		Tokens:        Token{},
		Subscriptions: Subscription{},
	}
}

func getInsertID(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
