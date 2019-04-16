package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var UserName = "root"  //user Name
var Password = "hs123456" //password
var DbName = "go_test"  //database Name

func main() {

	db, err := sql.Open("mysql", UserName+":"+Password+
		"@tcp(127.0.0.1:3306)/"+DbName)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	stmt, err := db.Prepare("CREATE TABLE person (id int NOT NULL AUTO_INCREMENT, first_name varchar(40), last_name varchar(40), PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Printf("Person Table successfully migrated....")
	}
}
