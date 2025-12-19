package main

// package post05

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type Userdata struct {
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

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// func exists(username string) int {
// 	username = strings.ToLower(username)
// 	db, err := openConnection()
// 	if err != nil {
// 		fmt.Println("exists", err)
// 		return -1
// 	}
// 	defer db.Close()

// 	rowsCount, err := checkIfEmptyDb()
// 	fmt.Println("exists: after rowsCount")

// 	if rowsCount == 0 {
// 		strErr := "exists: rowsCount === 0"

// 		if err != nil {
// 			strErr = err.Error()
// 		}
// 		fmt.Println(strErr)
// 		return 0
// 	}

// 	userID := -1
// 	statement := fmt.Sprintf("SELECT 'id' FROM 'users' where username='%s'", username)
// 	rows, err := db.Query(statement)
// 	for rows.Next() {
// 		var id int
// 		err = rows.Scan(&id)
// 		if err != nil {
// 			fmt.Println("Scan", err)
// 			return -1
// 		}
// 		userID = id
// 	}
// 	defer rows.Close()
// 	return userID
// }

// func AddUser(d Userdata) int {
// 	d.Username = strings.ToLower(d.Username)
// 	db, err := openConnection()
// 	if err != nil {
// 		fmt.Println(err)
// 		return -1
// 	}
// 	defer db.Close()
// 	rowsCount, err := isEmptyDb()

// 	fmt.Println("AddUser: rowsCount", rowsCount)

// 	if err != nil {
// 		fmt.Println("AddUser: error counting rows:", err)
// 		return -1
// 	}

// 	// if rowsCount == 0 {
// 	// 	return 0
// 	// }

// 	userID := exists(d.Username)
// 	fmt.Println("AddUser: userID:", userID)
// 	if userID > 0 {
// 		fmt.Println("AddUser: user already exists:", Username)
// 		return -1
// 	}
// 	fmt.Println("__after 1st exists__", d.Username)
// 	insertStatement := `insert into "users" ("username") values ($1)`
// 	// параметр $1 передаём в запрос аргументом db.Exec
// 	_, err = db.Exec(insertStatement, d.Username)
// 	if err != nil {
// 		fmt.Println(err)
// 		return -1
// 	}
// 	userID = exists(d.Username)
// 	if userID <= 0 {
// 		fmt.Println("AdUser: user was not added")
// 		// return userID
// 	}
// 	fmt.Println("AddUser: user added to table users, user id:", userID)
// 	insertStatement = `insert into "userdata" ("userid","name","surname","description") values ($1,$2,$3,$4)`
// 	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
// 	if err != nil {
// 		fmt.Println("db.Exec()", err)
// 		return -1
// 	}
// 	return userID
// }

// func DeleteUser(id int) error {
// 	db, err := openConnection()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
// 	// Провекра существования пользователя
// 	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
// 	rows, err := db.Query(statement)
// 	var username string
// 	for rows.Next() {
// 		err = rows.Scan(&username)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	defer rows.Close()
// 	if exists(username) != id {
// 		return fmt.Errorf("User with ID %d does not exist", id)
// 	}
// 	deleteStatement := `DELETE FROM "userdata" WHERE userid=$1`
// 	_, err = db.Exec(deleteStatement, id)
// 	if err != nil {
// 		return err
// 	}
// 	deleteStatement = `DELETE FROM "users" WHERE id=$1`
// 	_, err = db.Exec(deleteStatement, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
							"id","username","name","surname","description" 
						FROM "users","userdata" 
						WHERE users.id=userdata.userid`)
	// fmt.Println("ListUsers: before main query")
	if err != nil {
		fmt.Println("ListUsers: rows query error")
		return Data, err
	}
	// fmt.Println("ListUsers: after main query")

	for rows.Next() {
		var id int
		var username, name, surname, description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		// fmt.Println("ListUsers: in a cycle: temp Userdata: ", temp)
		Data = append(Data, temp)
	}

	if err != nil {
		return Data, err
	}

	defer rows.Close()
	return Data, nil
}

func checkIfEmptyDb() (bool, error) {
	db, err := openConnection()
	if err != nil {
		return true, err
	}
	defer db.Close()

	var countRows int
	err = db.QueryRow(`SELECT COUNT(*) FROM "users"`).Scan(&countRows)

	// fmt.Println("checkIfEmptyDb: users: err:", err)             // nil
	// fmt.Println("checkIfEmptyDb: users: countRows:", countRows) // 1

	if err != nil {
		return true, fmt.Errorf("error counting users:", err)
	}

	if countRows == 0 {
		fmt.Println("checkIfEmptyDb: users - empty")
		return true, nil
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM "userdata"`).Scan(&countRows)

	// fmt.Println("checkIfEmptyDb: userdata: err:", err)             // nil
	// fmt.Println("checkIfEmptyDb: userdata: countRows:", countRows) // 1

	if err != nil {
		return true, fmt.Errorf("error counting userdata:", err)
	}

	if countRows == 0 {
		fmt.Println("checkIfEmptyDb: userdata - empty")
		return true, nil
	}

	return false, nil
}

// func UpdateUser(d Userdata) error {
// 	db, err := openConnection()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	userID := exists(d.Username)
// 	if userID == -1 {
// 		return errors.New("User does not exist")
// 	}

// 	d.ID = userID
// 	updateStatement := `UPDATE "userdata" SET "name"=$1,"surname"=$2,"description"=$3 WHERE "userid"=$4`
// 	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// // генерирование случайной строки
func getString(length int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
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

func main() {
	Hostname = "localhost"
	Port = 5432
	Username = "postgres"
	Password = "mysecretpassword"
	Database = "postgres"
	isEmptyDb, err := checkIfEmptyDb()

	if err != nil {
		fmt.Println("func main: error check if empty db:", err)
		return
	}

	if isEmptyDb {
		fmt.Println("func main: db is empty, no ListUsers will be called")
	}

	if !isEmptyDb {
		data, err := ListUsers()

		if err != nil {
			fmt.Println("error using ListUsers:", err)
		}

		for _, v := range data {
			fmt.Println("func main - row:", v)
		}
	}

	SEED := time.Now().Unix()
	rand.New(rand.NewSource(SEED)) //  rand.Seed(SEED) устарело
	random_username := getString(5)

	fmt.Println("random_username:", random_username) // correct
	// isUserExists := exists(random_username)

	// t := Userdata{
	// 	Username:    random_username,
	// 	Name:        "Test Name post05",
	// 	Surname:     "Test Surname post05",
	// 	Description: "Test Description post05",
	// }

	// id := AddUser(t)

	// fmt.Println("__user added__")

	// if id == -1 {
	// 	fmt.Println("Error adding user", t.Username)
	// }

	// err = DeleteUser(id)
	// if err != nil {
	// 	fmt.Println("error using DeleteUser", err)
	// }
	// AddUser(t)
	// if id == -1 {
	// 	fmt.Println("error adding 2nd user", t.Username)
	// }
	// t = Userdata{
	// 	Username:    random_username,
	// 	Name:        "Test",
	// 	Surname:     "User 1",
	// 	Description: "this night not be me",
	// }
	// err = UpdateUser(t)
	// if err != nil {
	// 	fmt.Println("error using UpdateUser", err)
	// }

}
