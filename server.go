package main

import(
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/accounts", createAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", deleteAccount).Methods("DELETE")
	router.HandleFunc("/accounts", updateAccount).Methods("PUT")
	router.HandleFunc("/accounts", getAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}", getAccountByID).Methods("GET")
	
	log.Fatal(http.ListenAndServe("localhost:8000", router))
	
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

