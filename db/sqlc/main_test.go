package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DakshChawla/simplebank/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Println(err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Println(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
