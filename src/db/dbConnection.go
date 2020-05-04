package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "encoding/json"
	"github.com/tkanos/gonfig"
	"os"
	"log"
)
// Db for model to manipulate the data
type Db struct {}

type Configuration struct {
	Username string
	Password string
	Host string
	Database string
}

// CreateCon for setting up mysql connection
func (traceDb *Db) CreateCon() *sql.DB {
	dir, err := os.Getwd()
	
	if err != nil {
        log.Fatal(err)
	}

	
	configuration := Configuration{}
	err = gonfig.GetConf(dir + "/config.json", &configuration)
	
	if err != nil {
		panic(err)
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", configuration.Username, configuration.Password, configuration.Host, configuration.Database)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("data base error")
	} else {
		fmt.Println("data base connected successfully")
	}
	
	return db
}

func (traceDb *Db) Query(sql string, args ...interface{}) *sql.Rows{
	mydb := traceDb.CreateCon();
	result, err := mydb.Query(sql, args...)
	
	if err != nil {
		panic(err.Error())
	}

	result.Close()
	mydb.Close()
	
	return result
}

func (traceDb *Db) QueryRow(sql string, args ...interface{}) *sql.Row {
	mydb := traceDb.CreateCon();
	row := mydb.QueryRow(sql, args...)
	mydb.Close()
	return row
}


