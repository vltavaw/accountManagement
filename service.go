package main

import(
	"fmt"
	"net/http"
	"context"
	"time"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "encoding/json"
    "strconv"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Account struct {
	Id int
	Name string
	Amount int
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()

	var account Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)

	stmt, err := db.Prepare("INSERT into account SET Id=?, Name=?, Amount=?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(account.Id, account.Name, account.Amount)
	if err2 != nil {
		fmt.Println(err2)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account created successfully."))
}


func deleteAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	stmt, err := db.Prepare("DELETE FROM account WHERE id=?")
	if err != nil{
		fmt.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil{
		fmt.Println(err)
	}
	json, _ := json.MarshalIndent(result, "", "  ")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()

	var account Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)

	stmt, err := db.Prepare("UPDATE account SET Name=?, Amount=? WHERE id=?")
	if err != nil{
		fmt.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(account.Name, account.Amount, account.Id)
	if err != nil{
		fmt.Println(err)
	}

	cnt, _ := result.RowsAffected();
	if(cnt == 0) {
		w.Write([]byte("Could not find account with id: " + strconv.Itoa(account.Id)))
	} else { // check redis
		rdb := connectRedisClient()
		defer rdb.Close()
		params := mux.Vars(r)
		_, err := rdb.Get(ctx, params["id"]).Result()
		if(err == nil) {
			result := []Account{}
			result = append(result, account)
    		json, _ := json.MarshalIndent(result, "", "  ")
			rdb.Set(ctx, params["id"], json, time.Minute)
		} else if (err != redis.Nil) {
			fmt.Println(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account updated successfully."))
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM account")
	if err != nil{
		fmt.Println(err)
	}
	defer rows.Close()

	json, err2 := printJsonResults(rows)
	if err2 != nil{
		fmt.Println(err2)
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func getAccountByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	rdb := connectRedisClient()
	defer rdb.Close()
	val, err := rdb.Get(ctx, params["id"]).Result()

	if(err == redis.Nil) { // not cached in redis
		db := connectDB()
		defer db.Close()
		rows, err := db.Query("SELECT * FROM account WHERE Id=?", id)
		if err != nil {
			fmt.Println(err)
		}
    	defer rows.Close()

   		if rows.Next() == false {
    		w.WriteHeader(http.StatusOK)
    		w.Write([]byte("Could not find account with id: " + strconv.Itoa(id)))
    		return
    	}

    	result := []Account{}
		account := Account{}
		err2 := rows.Scan(&account.Id, &account.Name, &account.Amount)
    	if err2 != nil {
        	fmt.Println(err)
    	}
    	result = append(result, account)
    	json, _ := json.MarshalIndent(result, "", "  ")
	
		w.WriteHeader(http.StatusOK)
		w.Write(json)

		// write into redis
		rdb.Set(ctx, params["id"], json, time.Minute)
		return
	} else if (err != nil) {
		fmt.Println(err)
		return
	} else { // read from redis
		account := Account{}
		result := []Account{}
		json.Unmarshal([]byte(val), &account) // formatting stuff
		result = append(result, account)
		json, _ := json.MarshalIndent(result, "", "  ")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
		return
	}	
}