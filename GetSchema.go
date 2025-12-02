package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	// Пакет lib/pq, который является интерфейсом к базе данных PostgreSQL, не используется непосредственно кодом. Следовательно, вам необходимо импортировать пакет lib/pq с помощью _, чтобы компилятор Go не выдавал сообщение об ошибке, связанное с импортом пакета, а не с его применением. Обычно это не нужно, но этот вид импорта обычно вызван тем, что импортированный пакет имеет побочные эффекты, такие как регистрация себя в качестве обработчика базы данных для пакета sql
	_ "github.com/lib/pq"
)

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

	conn := fmt.SPrintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Open():", err)
		return
	}
	defer db.Close()

	// получить все БД. Текущий результат запроса SELECT сохраняется в переменной rows, которая служит курсором. Вы не получите все результаты из базы данных сразу (так как запрос может возвращать миллионы записей), но получите их одну за другой — в этом весь смысл использования курсора.
	rows, err := db.Query(`SELECT "datname" FROM "pg_database" WHERE datistemplate=false`)
	if err != nil {
		fmt.Println("Query:", err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan"), err
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()
}
