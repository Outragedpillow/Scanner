package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"strings"
	// "Scanner/sqlite"
	"Scanner/api"
)

const (
	RESIDENT_FILEPATH string = "./data/residents.txt"
	COMPUTER_FILEPATH string = "./data/computers.txt"
)

func ProcessScan(db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	// we are passing a db pointer in a loop here,
	// so lets make sure it's the input we want
	for scanner.Scan() {
		input := scanner.Text()
		if len(input) == 0 {
			continue
		} else if len(input) == 8 {
			FindComputer(db, input)
		}
	}
}

func FindComputer(db *sql.DB, serial string) {
	var computer api.Computer

	sqlStatement, prepErr := db.Prepare("SELECT serial, tag_number, is_issued, signed_out_by, signed_out_to, time_issued, time_returned FROM computers WHERE serial = ?")
	if prepErr != nil {
		fmt.Println("Error: Prepare")
	}

	defer sqlStatement.Close()

	row := sqlStatement.QueryRow(serial)

	err := row.Scan(&computer.Serial, &computer.Tag_number, &computer.Signed_out_by, &computer.Signed_out_to, &computer.Time_issued, &computer.Time_returned)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No row.")
		}
	}

	fmt.Println(computer.Tag_number)
}

func ReadComputersIntoDb(db *sql.DB) error {
	file, openErr := os.Open(COMPUTER_FILEPATH)
	if openErr != nil {
		return openErr
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println("Scanning")
		if len(scanner.Text()) == 0 {
			continue
		}
		comp := parseInput(scanner.Text())
		insertErr := insertComputerData(db, comp)
		if insertErr != nil {
			fmt.Println(insertErr)
		}
	}

	return nil
}

func parseInput(info string) []string {

	words := strings.Split(info, " ")

	switch len(words) {
	case 3:
		words[0] = strings.Trim(words[0], ",")
		return words

	case 2:
		serial := words[0]
		if len(serial) != 20 {
			return nil
		}
		index := strings.Index(serial, "R")
		if index != -1 {
			serial = serial[index:]
			words[0] = serial
			return words
		}
	}
	return nil
}

func insertComputerData(db *sql.DB, info []string) error {
	if len(info) == 2 {
		sqlStatement, prepErr := db.Prepare("INSERT INTO computers (serial, tag_number, is_issued) VALUES (?, ?, ?)")
		if prepErr != nil {
			return prepErr
		}

		defer sqlStatement.Close()
		// We are inserting 0 for is_issued?
		_, execErr := sqlStatement.Exec(info[0], info[1], 0)
		if execErr != nil {
			return execErr
		}

		return nil
	}

	return errors.New("Invalid input length.")
}
func ReadResidentsIntoDb(db *sql.DB) {
	file, openErr := os.Open(RESIDENT_FILEPATH)
	if openErr != nil {
		fmt.Println("Error", openErr)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res := parseInput(scanner.Text())
		insertErr := insertResidentData(db, res)
		if insertErr != nil {
			fmt.Println("Error", insertErr)
			return
		}
	}
}

func insertResidentData(db *sql.DB, info []string) error {
	mdoc, convErr := strconv.Atoi(info[2])
	if convErr != nil {
		return convErr
	}

	name_of := info[0] + " " + info[1]
	sqlStatement, prepErr := db.Prepare("INSERT INTO residents (name_of, mdoc) values (?, ?)")
	if prepErr != nil {
		return prepErr
	}

	defer sqlStatement.Close()

	_, execErr := sqlStatement.Exec(name_of, mdoc)
	if execErr != nil {
		return execErr
	}

	return nil
}
