package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func CreateDataBase(user string, pw string, domain string, port string, dbName string) *DataBase {
	db := InitDatabase(user, pw, domain, port, dbName)
	CreateTables(db)
	preparedStatements := make(map[string]*sql.Stmt)
	stmt, err := db.Prepare("INSERT INTO tea (name, link, data) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE name=?, data=?")
	if err != nil {
		log.Fatal(err)
	}

	preparedStatements["add"] = stmt

	return &DataBase{
		db,
		preparedStatements,
	}
}

type DataBase struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func (d *DataBase) AddTea(name string, link string, data string) {
	// TODO upsert teas

	_, err := d.statements["add"].Exec(name, link, data, name, data)
	if err != nil {
		fmt.Printf("error in AddTea: err: %s", err)
		fmt.Println("")
	}
}

func (d *DataBase) Prepare(statement string) *sql.Stmt {
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}

	return stmt
}

func CreateTables(db *sql.DB) {
	query := "CREATE TABLE IF NOT EXISTS tea (id serial, name VARCHAR(255), link VARCHAR(255), data TEXT, UNIQUE (link))"
	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func InitDatabase(user string, pw string, domain string, port string, dbName string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pw, domain, port, dbName))

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

	if err != nil {
		// db may already EXISTSist

		_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err.Error())
		}

	} else {
		_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
		if err != nil {
			panic(err)
		}
	}

	return db
}
