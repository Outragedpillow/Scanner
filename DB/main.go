package main

import (
	"Scanner/sqlite"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const PORT = ":1234"

func main() {

	err := checkDbExists()
	if err != nil {
		fmt.Println("Failed to remove existing db")
	}
	// Here we create database, and get pointer to it
	db, tableErr := sqlite.CreateTables()
	if tableErr != nil {
		fmt.Println("Failed to Create tables.")
		return
	}
	fmt.Println("Read files")
	ReadResidentsIntoDb(db)

	fmt.Println("Read files 2")
	ReadComputersIntoDb(db)

	fmt.Println("Read files 3")
	var wg sync.WaitGroup
	wg.Add(1)
	// Before we start our HTTP server, we register our handler function
	// This is an example function of the handler for the /computers route
	http.HandleFunc("/computers", handleComputerInsertion(db))
	go func() {
		wg.Done()
		err := http.ListenAndServe(PORT, nil)
		if err != nil {
			fmt.Println("Failed to start server")
		}
		fmt.Println("Server Started")
	}()

	wg.Wait()

	fmt.Println("ScanScanning")
	ProcessScan(db)

}

func checkDbExists() error {
	file, err := os.Stat("./Storage.db")
	if err == nil {
		rmErr := os.Remove(file.Name())
		if rmErr != nil {
			return rmErr
		}
	}
	return nil
}

// This function is called by our HTTP server, and is responsible for
// inserting the computer data into the database as well as responding to the
// HTTP request with a success message.

func handleComputerInsertion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Our initial read of the HTTP request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Request body error", http.StatusInternalServerError)
			return
		}

		data := string(body)
		if len(data) == 0 {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		// our newly defined parseInput function returns nil if invalid
		comp := parseInput(data)
		if comp == nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		// Else we insert the computer data into the database
		if err := insertComputerData(db, comp); err != nil {
			http.Error(w, "Error inserting data into the database", http.StatusInternalServerError)
			return
		}

		// Respond with a success message :D
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Data inserted successfully")
	}
}
