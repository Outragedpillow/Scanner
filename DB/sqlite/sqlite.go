package sqlite

import (
	"fmt"
	// "strconv"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CreateTables() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "Storage.db")
	if err != nil {
		fmt.Println("Open Error", err)
		return nil, err
	}

	_, err = database.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		fmt.Println("PRAGMA Error", err)
		return nil, err
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS residents (
    name_of TEXT NOT NULL,
    mdoc INTEGER PRIMARY KEY NOT NULL);
  `)
	if err != nil {
		fmt.Println("Create Table Error: Residents.", err)
		return nil, err
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS admin (
    name_of TEXT PRIMARY KEY NOT NULL);
  `)
	if err != nil {
		fmt.Println("Create Table Error: Admin.", err)
		return nil, err
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS computers (
    serial TEXT PRIMARY KEY NOT NULL,
    tag_number INT NOT NULL,
    is_issued INT NOT NULL,
    signed_out_by INTEGER,
    signed_out_to INTEGER,
    time_issued TIMESTAMP,
    time_returned TIMESTAMP,
    FOREIGN KEY(signed_out_to) REFERENCES residents(mdoc)
    );
  `)
	if err != nil {
		fmt.Println("Create Table Error: Computers.", err)
		return nil, err
	}

	triggerStatement := fmt.Sprintf(`
      CREATE TRIGGER check_issued 
      BEFORE INSERT ON computers 
      BEGIN 
        SELECT CASE 
            WHEN NEW.is_issued = 1 AND (NEW.signed_out_by IS NULL OR NEW.signed_out_to IS NULL) THEN
              RAISE (ABORT, 'CANNOT ISSUE WITHOUT VALUES FOR signed_out_to and signed_out_by.')
        END;
      END;
    `)

	_, err = database.Exec(triggerStatement)
	if err != nil {
		fmt.Println("Trigger creation error:", err)
	}

	return database, nil
}
