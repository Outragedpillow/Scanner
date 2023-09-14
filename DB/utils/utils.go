package utils

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "strings"
  "bufio"
  "os"
  // "Scanner/sqlite"
  "Scanner/structs"
  "strconv"
)

type Scans struct {
  Scan1 interface{}
  Scan2 interface{}
}

var scanned Scans = Scans{}

func ProcessScan(db *sql.DB) {

  scanner := bufio.NewScanner(os.Stdin);
  for { 
    scanned = Scans{};
    for i:=0; i<2; i++ {
      if scanner.Scan() {
        input := scanner.Text();
        if len(input) > 2 {    
          switch input[0:2] {
          case "1s":
            findCompErr := findComputer(db, input);
            if findCompErr != nil {
              fmt.Println("Break", findCompErr);
              break;
            }
          default:
            resMdoc, convErr := strconv.Atoi(input);
            if convErr != nil {
              fmt.Println("Error: Convert ", convErr);
            }
            findResErr := findResident(db, resMdoc);
            if findResErr != nil {
              fmt.Println("Break", findResErr);
              break;
            }
          }
        } 
      }
    }

    switch scanned.Scan1.(type) {
      case structs.Computer:
        if res, ok := scanned.Scan2.(structs.Resident); ok {
          if comp, ok := scanned.Scan1.(structs.Computer); ok {
            fmt.Println("Setup for Api");
            data := fmt.Sprintf("\nResident name: %s, MDOC: %d, Computer s/n: %s, Computer tag number: %d", res.Name_of, res.Mdoc, comp.Serial, comp.Tag_number);
            WriteComputerLogs(data);
            fmt.Println("Signed out to:", res.Name_of, "| Computer number:", comp.Tag_number);  
          } else {
          fmt.Println("Error: Wrong combination of scans");
          }
        }

      default:
        fmt.Println("Error: Default wrong combination of scans");
      }
  }
} 

func findComputer(db *sql.DB, serial string) error {
   var computer structs.Computer;

  fmt.Println("finding Computer")

  // Split input to get just serial number
  index := strings.Index(serial, "R");
  if index != -1 {
    serial = serial[index:];
  }

  // Prepare statement for select fields of table computer in sqlite and join foreign keys where serial is input from scanner after being sliced
  sqlStatement, prepErr := db.Prepare("SELECT c.serial, c.tag_number, c.is_issued, a.name_of_a, r.name_of_r, c.time_issued, c.time_returned FROM computers AS c LEFT JOIN admin AS a ON c.signed_out_by = a.name_of_a LEFT JOIN residents AS r ON c.signed_out_to = r.mdoc WHERE serial = ?");
  if prepErr != nil {
    fmt.Println("Error: Prepare", prepErr)
    return prepErr;
  }

  defer sqlStatement.Close();

  // Execute prepared statement
  row := sqlStatement.QueryRow(serial);
  
  // Iterate over fields of selected row and asssign the values to computer struct 
  rowErr := row.Scan(&computer.Serial, &computer.Tag_number, &computer.Is_issued, &computer.Signed_out_by.Name_of, &computer.Signed_out_to.Mdoc, &computer.Time_issued, &computer.Time_returned);
  if rowErr != nil {
    if rowErr == sql.ErrNoRows {
      fmt.Println("No row. ", rowErr);
      return rowErr;
    }
    fmt.Println(rowErr);
  }
  
  fmt.Println(computer.Tag_number);
  
  // tagNumberStr := strconv.Itoa(computer.Tag_number);

  // Set Scan1 as computer for later validation of both scans
  scanned.Scan1 = computer;
  return nil;
}

func findResident(db *sql.DB, mdoc int) error {
  var resident structs.Resident;

  fmt.Println("Finding Resident");

  sqlStatement, prepErr := db.Prepare("SELECT mdoc, name_of_r FROM residents WHERE mdoc = ?");
  if prepErr != nil {
    fmt.Println("Error: Prepare ", prepErr);
    return prepErr;
  }

  defer sqlStatement.Close()

  row := sqlStatement.QueryRow(mdoc);

  rowErr := row.Scan(&resident.Mdoc, &resident.Name_of)
  if rowErr != nil {
    if rowErr == sql.ErrNoRows {
      fmt.Println("Error: No row. ", rowErr);
      return rowErr;
    } else {
      fmt.Println("Error: Row scan error: ", rowErr);
    }
  }

  fmt.Println(resident.Name_of)

  scanned.Scan2 = resident;
  return nil;
}


