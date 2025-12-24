package post05

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"

	_ "github.com/lib/pq"
)

type UserData struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

var (
	Hostname = ""
	Port     = 5432
	Username = ""
	Password = ""
	Database = ""
)

var isEmptyDb bool
var MIN = 0
var MAX = 26

func openDbCon() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func FindUserId(username string) int {
	username = strings.ToLower(username)
	db, err := openDbCon()
	if err != nil {
		fmt.Println("findUserID: could not connect db. error: ", err)
		return -1
	}
	defer db.Close()

	userId := -1
	statement := fmt.Sprintf("SELECT id FROM users WHERE username='%s'", username)
	rows, err := db.Query(statement)

	if err != nil {
		fmt.Println("findUserId: error on select by id:", err)
		return -1
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("findUserId: Scan error:", err)
			return -1
		}
		userId = id
	}
	defer rows.Close()
	return userId
}

func AddUser(udata UserData) int {
	udata.Username = strings.ToLower(udata.Username)
	db, err := openDbCon()
	if err != nil {
		fmt.Println("AddUser: error open db connection:", err)
		return -1
	}
	defer db.Close()

	insertStatement := `INSERT INTO "users" ("username") VALUES ($1)`
	_, err = db.Exec(insertStatement, udata.Username)
	if err != nil {
		fmt.Println("AddUser: error db.Exec:", err)
		return -1
	}

	newUserId := FindUserId(udata.Username)
	if newUserId <= 0 {
		fmt.Println("AdUser: user was not added to table users")
		return -1
	}
	fmt.Println("AddUser: user added to table users. Id:", newUserId)
	insertStatement = `INSERT INTO "userdata" ("userid","name","surname","description") VALUES ($1,$2,$3,$4)`
	_, err = db.Exec(insertStatement, newUserId, udata.Name, udata.Surname, udata.Description)
	if err != nil {
		fmt.Println("AdUser: user was not added to table userdata. Error:", err)
		return -1
	}
	fmt.Println("AddUser: user added to userdata. newUserId:", newUserId)
	return newUserId
}

func DeleteUser(id int) error {
	db, err := openDbCon()
	if err != nil {
		return err
	}
	defer db.Close()

	q := `DELETE FROM "userdata" WHERE userid=$1`
	_, err = db.Exec(q, id)
	if err != nil {
		return err
	}

	q = `DELETE FROM "users" WHERE id=$1`
	_, err = db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func ListUsers() ([]UserData, error) {
	Data := []UserData{}
	db, err := openDbCon()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
							"id","username","name","surname","description" 
						FROM "users","userdata" 
						WHERE users.id=userdata.userid`)
	if err != nil {
		fmt.Println("ListUsers: rows query error")
		return Data, err
	}

	for rows.Next() {
		var id int
		var username, name, surname, description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := UserData{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
	}

	if err != nil {
		return Data, err
	}

	defer rows.Close()
	return Data, nil
}

func CheckIfEmptyDb() (bool, error) {
	db, err := openDbCon()
	if err != nil {
		return true, err
	}
	defer db.Close()

	var countRows int
	err = db.QueryRow(`SELECT COUNT(*) FROM "users"`).Scan(&countRows)

	if err != nil {
		return true, fmt.Errorf("error counting users:", err)
	}

	if countRows == 0 {
		fmt.Println("CheckIfEmptyDb: users - empty")
		return true, nil
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM "userdata"`).Scan(&countRows)

	if err != nil {
		return true, fmt.Errorf("error counting userdata:", err)
	}

	if countRows == 0 {
		fmt.Println("CheckIfEmptyDb: userdata - empty")
		return true, nil
	}

	return false, nil
}

func UpdateUser(d UserData) error {
	db, err := openDbCon()
	if err != nil {
		return err
	}
	defer db.Close()

	updateStatement := `UPDATE "userdata" SET "name"=$1,"surname"=$2,"description"=$3 WHERE "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// Генерирование случайной строки
func GetString(length int64) string {
	startChar := "A"
	var i int64 = 1
	temp := ""

	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar

		if i == length {
			break
		}
		i++
	}
	return temp
}
