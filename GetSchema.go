package main

// Утилита проверяет возможность подключиться к базе данных PostgreSQL и получить список доступных баз данных и таблиц в этой базе данных и public схеме — вся информация о подключении предоставляется в качестве аргументов командной строки

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func findTables() {
	query := `SELECT table_name FROM information_schema.tables WHERE atble_schema='public' ORDER BY table_name`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Query:", err)
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("+T", name)
	}
	defer rows.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) != 6 {
		fmt.Println("Please provide: hostname port username password dbname")
		return
	}
	port, err := strconv.ParseInt(arguments[2], 10, 64)
	if err != nil {
		fmt.Println("Not a valid port number:", err)
	}
	host := arguments[1]
	user := arguments[3]
	pass := arguments[4]
	database := arguments[5]

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Open():", err)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "datname" FROM "pg_database" WHERE datistemplate=false`)
	if err != nil {
		fmt.Println("Query:", err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()
}
