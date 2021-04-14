package main

import(
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
)

type Account struct {
	Id int
	Name string
	Amount int
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/accounts", createAccount).Methods("POST")
	router.HandleFunc("/accounts", deleteAccount).Methods("DELETE")
	router.HandleFunc("/accounts", updateAccount).Methods("PUT")
	router.HandleFunc("/accounts", getAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}", getAccountByID).Methods("GET")
	
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:mysqlpassword233@tcp(127.0.0.1:3306)/accountdb")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Println("Successfully connected.")
	return db
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()

	var account Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)

	stmt, err := db.Prepare("INSERT into account SET Name=?,Amount=?")
	if err != nil {
		fmt.Println(err)
	}

	_, err2 := stmt.Exec(account.Name, account.Amount)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()

	rows, err := db.Query("SELECT * FROM account")
	if err != nil{
		fmt.Println(err)
	}

	results := []Account{}
	account := Account{}

	for rows.Next() {
        e := rows.Scan(&account.Id, &account.Name, &account.Amount)
        if e != nil {
            fmt.Println(err)
        }
        results = append(results, account)
    }

    rows.Close()
	json, _ := json.Marshal(results)
	w.Write(json)
}

func getAccountByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
