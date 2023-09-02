package utils

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  // "strings"
  "bufio"
  "os"
  // "Scanner/sqlite"
  "Scanner/api"
)

func ProcessScan(db *sql.DB) {
  scanner := bufio.NewScanner(os.Stdin);
  for scanner.Scan() {
    input := scanner.Text();
    checkScanType(db, input);
  }
} 

func checkScanType(db *sql.DB, input string) {
  if len(input) == 8 {
    FindComputer(db, input);
  }
}


func FindComputer(db *sql.DB, serial string) {
  var computer api.Computer;

  sqlStatement, prepErr := db.Prepare("SELECT serial, tag_number, is_issued, signed_out_by, signed_out_to, time_issued, time_returned FROM computers WHERE serial = ?");
  if prepErr != nil {
    fmt.Println("Error: Prepare")
  }

  defer sqlStatement.Close();

  row := sqlStatement.QueryRow(serial);

  err := row.Scan(&computer.Serial, &computer.Tag_number, &computer.Signed_out_by, &computer.Signed_out_to, &computer.Time_issued, &computer.Time_returned);
  if err != nil {
    if err == sql.ErrNoRows {
      fmt.Println("No row.")
    }
  }
  
  fmt.Println(computer.Tag_number);
}
