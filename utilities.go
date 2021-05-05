package main

import(
	"fmt"
	"database/sql"
    "encoding/json"
    "github.com/go-redis/redis/v8"
)

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:mysqlpassword233@tcp(127.0.0.1:3306)/accountdb")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Println("Successfully connected to database.")
	return db
}

func connectRedisClient() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "127.0.0.1:6379",
        Password: "",
        DB:       0,
    })

    return rdb
}

func connectRedis() int {
	return 1
}

func printJsonResults(rows *sql.Rows) ([]byte, error) {
	results := []Account{}
	account := Account{}

	for rows.Next() {
        e := rows.Scan(&account.Id, &account.Name, &account.Amount)
        if e != nil {
            return nil, e
        }
        results = append(results, account)
    }
  	
	json, _ := json.MarshalIndent(results, "", "  ")
	return json, nil
}