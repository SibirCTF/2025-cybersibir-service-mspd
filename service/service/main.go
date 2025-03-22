package main

import (
	"log"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/router"
)

func main() {
	DB, err := db.Init_db()
	if err != nil {
		log.Fatal(DB, err)
	}
	router.Setup_routes()
	err = http.ListenAndServe("0.0.0.0:1015", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
