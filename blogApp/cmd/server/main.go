package main

import (
	dbconfig "blogApp/internal/config"
	handler "blogApp/internal/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//Load env variables
	godotenv.Load()
	db_name := os.Getenv("DB_NAME")

	// database connection
	client, err := dbconfig.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(db_name)

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/signup", handler.SingnUp(db)).Methods("POST")
	router.HandleFunc("/{users}", handler.GetUsers(db)).Methods("GET")

	fmt.Println("Server is running at port 9000")
	http.ListenAndServe(":9000", router)
}
