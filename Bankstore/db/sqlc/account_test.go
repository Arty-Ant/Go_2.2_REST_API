/*
	Pattern: <filename>_test.go
*/

package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://app_user:pswd@localhost:5432/bankdb?sslmode=disable"
)

var ctx = context.Background()

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	defer conn.Close(ctx)

	testQueries = New(conn)

	os.Exit(m.Run())
}

func TestCreateAccount(t *testing.T) {
	fmt.Println("Run test")
}
