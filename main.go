package main

import (
	"apiBooks/api"
	"apiBooks/database"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	chOff := make(chan os.Signal)
	signal.Notify(chOff, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	host, dbPort, user, password, dbname := "localhost", "5432", "Lilit", "testtest", "apiBooks"
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbname)

	conn, err := sql.Open("postgres", psqlConn)
	if err != nil {
		log.Fatal("Start database:", err)
	}

	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal("Ping database:", err)
	}

	log.Println("Ping database OK")

	db := database.NewStorage(conn)

	port := ":8080"
	router := mux.NewRouter()

	server := api.NewServer(port, router, db)
	log.Println("Start server OK")

	go func() {
		err = server.Start()
		if err != nil {
			log.Fatal("Start server:", err)
		}
	}()

	<-chOff
}
