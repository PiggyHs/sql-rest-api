package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id         int
	First_Name string
	Last_Name  string
}

var userName = "root"  //user Name
var password = "hs123456" //password
var dbName = "go_test"  //database Name

func main() {
	db, err := sql.Open("mysql", userName+":"+password+
		"@tcp(127.0.0.1:3306)/"+dbName)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	engine := gin.Default()

	//select person by id
	engine.GET("/person/:id", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id") //解析
		row := db.QueryRow("select id, first_name, last_name from person where id = ?;", id)
		err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	engine.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("select id, first_name, last_name from person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})

	engine.POST("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// PUT - update a person details
	engine.PUT("/person/:id", func(c *gin.Context) {
		var buffer bytes.Buffer
		cid := c.Param("id")
		id, err := strconv.Atoi(cid)
		if err != nil {
			panic(err)
		}
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources
	engine.DELETE("/person/:id", func(c *gin.Context) {
		cid := c.Param("id")
		id, err := strconv.Atoi(cid)
		if err != nil {
			panic(err)
		}
		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})

	engine.Run(":3000")
}