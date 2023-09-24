package structs

import (
	"database/sql"
	"fmt"
)

type Database struct {
  Conn *sql.DB
}

type Resident struct {
  Mdoc int `json:"mdoc"`
  Name_of string `json:"name_of"`
  Has_computer bool `json:"has_computer"`
  Issued_computer string `json:"Issued_computer"`
}

type Computer struct {
  Serial string `json:"serial"`
  Tag_number int `json:"tag_number"`
  Is_issued bool `json:"is_issued"`
  Signed_out_to Resident `json:"signed_out_to"`
  Time_issued string `json:"time_issued"`
  Time_returned string `json:"time_returned"`
}

func (db *Database) Close() error {
  if db.Conn != nil {
    return db.Conn.Close();
  }

  return nil;
}

func (database *Database) Open(name string) error {
  db, openErr := sql.Open("sqlite3", name);
  if openErr != nil {
    fmt.Println("Error: Open Storage.db ", openErr);
    return openErr;
  }

  database.Conn = db;
  return nil;
}

func (database *Database) CreateTables() error {
  _, err := database.Conn.Exec(`PRAGMA foreign_keys = ON;`);
  if err != nil {
    fmt.Println("PRAGMA Error", err);
    return err;
  }

  _, err = database.Conn.Exec(`CREATE TABLE IF NOT EXISTS residents (
    name_of_r TEXT NOT NULL,
    mdoc INTEGER PRIMARY KEY NOT NULL,
    has_computer INTEGER NOT NULL,
    issued_computer TEXT
    );
  `);
  if err != nil {
    fmt.Println("Create Table Error: Residents.", err);
    return err;
  }

  _, err = database.Conn.Exec(`CREATE TABLE IF NOT EXISTS computers (
    serial TEXT PRIMARY KEY NOT NULL,
    tag_number INTEGER NOT NULL,
    is_issued INTEGER NOT NULL,
    signed_out_to INTEGER,
    time_issued TIMESTAMP,
    time_returned TIMESTAMP,
    FOREIGN KEY(signed_out_to) REFERENCES residents(mdoc)
    );
  `);
  if err != nil {
    fmt.Println("Create Table Error: Computers.", err);
    return err;
  }
  
  
  triggerStatement := fmt.Sprintf(`
      CREATE TRIGGER check_issued 
      BEFORE INSERT ON computers 
      BEGIN 
        SELECT CASE 
            WHEN NEW.is_issued = 1 AND (NEW.signed_out_to IS NULL) THEN
              RAISE (ABORT, 'CANNOT ISSUE WITHOUT VALUES FOR signed_out_to.')
        END;
      END;
    `);

  _, err = database.Conn.Exec(triggerStatement);
  if err != nil {
    fmt.Println("Trigger creation error:", err)
  }

  return nil;
}

func (database *Database) IsSignedout(serial string) (int, error) {
  var is_issued int

  sqlStatement, prepErr := database.Conn.Prepare("SELECT is_issued FROM computers WHERE serial = ?")
  if prepErr != nil {
    return -1, fmt.Errorf("Error: Prepare statement, %v", prepErr)
  }

  defer sqlStatement.Close();

  queryErr := sqlStatement.QueryRow(serial).Scan(&is_issued)
  if queryErr != nil {
      if queryErr == sql.ErrNoRows {
        return -1, fmt.Errorf("Error: No result found, %v", queryErr)
      }
      return -1, fmt.Errorf("Error: Query, %v", queryErr)
  }

  return is_issued, nil
}

func (database *Database) IsSignedoutTo(serial string) (Resident, error) {
  var resident Resident;
  
  sqlStatement, prepErr := database.Conn.Prepare("SELECT r.mdoc, r.name_of_r FROM computers AS c LEFT JOIN residents AS r ON c.signed_out_to = r.mdoc WHERE serial = ?");
  if prepErr != nil {
    return resident, fmt.Errorf("Error: Prepare statement, %v", prepErr);
  }

  defer sqlStatement.Close();

  queryErr := sqlStatement.QueryRow(serial).Scan(&resident.Mdoc, &resident.Name_of);
  if queryErr != nil {
    if queryErr == sql.ErrNoRows {
      return resident, fmt.Errorf("Error: No result found, %v", queryErr);
    } else {
      return resident, fmt.Errorf("Error: Query, %v", queryErr);
    }
  }

  return resident, nil;
}

func (database *Database) HasComputer(mdoc int) (int, error) {
  var has_computer int;

  sqlStatement, prepErr := database.Conn.Prepare("SELECT has_computer FROM residents WHERE mdoc = ?");
  if prepErr != nil {
    return -1, fmt.Errorf("Error: Prepare statement, %v", prepErr);
  }

  defer sqlStatement.Close();

  queryErr := sqlStatement.QueryRow(mdoc).Scan(&has_computer);
  if queryErr != nil {
    if queryErr == sql.ErrNoRows {
      return -1, fmt.Errorf("Error: Query no rows, %v", queryErr);
    } else {
      return -1, fmt.Errorf("Error: Query statement, %v", queryErr);
    }
  }

  return has_computer, nil;
}

func (database *Database) HasComputerNumber(mdoc int) (Computer, error) {
  var computer Computer;

  sqlStatement, prepErr := database.Conn.Prepare("SELECT serial FROM computers WHERE signed_out_to = ?");
  if prepErr != nil {
    return computer, fmt.Errorf("Error: Prepare statement, %v", prepErr);
  }

  defer sqlStatement.Close();

  queryErr := sqlStatement.QueryRow(mdoc).Scan(&computer.Serial);
  if queryErr != nil {
    if queryErr == sql.ErrNoRows {
      return computer, fmt.Errorf("Error: Scan no rows, %v", queryErr);
    } else {
      return computer, fmt.Errorf("Error: Query statement, %v", queryErr);
    }
  }

  return computer, nil;
}

func ResidentIsEmpty(s Resident) bool {
  if s.Mdoc == 0 || s.Name_of == "" {
    return true;
  }

  return false;
}

func ComputerIsEmpty(c Computer) bool {
  if c.Serial == "" || c.Tag_number == 0 {
    return true;
  }

  return false;
}
