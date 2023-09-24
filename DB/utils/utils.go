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
  "time"
)

type Scans struct {
  Scan1 interface{}
  Scan2 interface{}
}

var scanned Scans = Scans{}
var CurrentSignOuts []string

func ProcessScan(database *structs.Database) {

  db := database.Conn;
  scanner := bufio.NewScanner(os.Stdin);
  for { 
    scanned = Scans{};
    for i:=0; i<2; i++ {
      if scanner.Scan() {
        input := scanner.Text();
        if len(input) > 2 {    
          switch input[0:2] {
          case "1s":
            // if computer is found will be assigned to scanned.Scan1 
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
            // if resident is found will be assigned to scanned.Scan2 
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
            // Add checks before continuing
            pass := checkComputerStatus(database, comp.Serial, res.Mdoc);
            switch pass {
            case 0:
              fmt.Println("Error: Computer can only be signed in by same person who signed it out.");
            case 1:
              fmt.Println("Setup for Api");
              updateErr := updateDbWithScans(db, &res, &comp);
              if updateErr != nil {
                fmt.Println("Failed to update");
              }
              data := fmt.Sprintf("Resident name: %s, MDOC: %d, Computer s/n: %s, Computer tag number: %d, Time issued: %s, Time returned: %s\n", res.Name_of, res.Mdoc, comp.Serial, comp.Tag_number, comp.Time_issued, comp.Time_returned);
              WriteComputerLogs(data, "history");
              deleteErr := DeleteSignedout();
              if deleteErr != nil {
                break;
              }
              CurrentSignOuts = updateCurrentSignOuts(CurrentSignOuts, data); 
              for _, v := range CurrentSignOuts {
                WriteComputerLogs(v, "signedout");
              }
            default:
              fmt.Println("Error: Default pass")
            }
          } else {
          fmt.Println("Error: Scan1 wrong combination of scans");
          }
        } else {
          fmt.Println("Error: Scan2 wrong combination of scans");
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
  sqlStatement, prepErr := db.Prepare("SELECT serial, tag_number, is_issued FROM computers WHERE serial = ?");
  if prepErr != nil {
    fmt.Println("Error: Prepare", prepErr);
    return prepErr;
  }

  defer sqlStatement.Close();

  // Execute prepared statement
  row := sqlStatement.QueryRow(serial);
  
  // Iterate over fields of selected row and asssign the values to computer struct 
  rowErr := row.Scan(&computer.Serial, &computer.Tag_number, &computer.Is_issued);
  if rowErr != nil {
    if rowErr == sql.ErrNoRows {
      fmt.Println("No row. ", rowErr);
      return rowErr;
    }
    fmt.Println("Error: Scan into computer", rowErr);
  }
  
  fmt.Println(computer.Tag_number);
  
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

func updateDbWithScans(db *sql.DB, res *structs.Resident, comp *structs.Computer) error {
  currentTime := time.Now();
  formattedTime := currentTime.Format("2006-01-02 15:04:05");


  err := func() error {

    if comp.Is_issued {
      fmt.Println("SIGNING OUT");
      comp.Is_issued = false;
      comp.Time_returned = formattedTime;

      sqlStatement, prepErr := db.Prepare("UPDATE computers SET is_issued = 0, signed_out_to = NULL, time_returned = ? WHERE serial = ?");
      if prepErr != nil {
        fmt.Println("Error: Prepare is_issued true ", prepErr);
        return prepErr;
      }

      defer sqlStatement.Close();

      _, execErr := sqlStatement.Exec(formattedTime, comp.Serial);
      if execErr != nil {
        fmt.Println("Error: Exec is_issued true ", execErr);
        return execErr;
      }

      return nil;

  } else if !comp.Is_issued {
    fmt.Println("SIGNING IN");
    comp.Is_issued = true;
    comp.Signed_out_to = *res;
    comp.Time_issued = formattedTime;

    sqlStatement, prepErr := db.Prepare("UPDATE computers SET is_issued = 1, Signed_out_to = ?, time_issued = ? WHERE serial = ?");
    if prepErr != nil {
      fmt.Println("Error: Prepare is_issued false", prepErr);
      return prepErr;
    }

    defer sqlStatement.Close();
    
    _, execErr := sqlStatement.Exec(res.Mdoc, formattedTime, comp.Serial);
    if execErr != nil {
      fmt.Println("Error: Exec is_issued false", execErr);
      return execErr;
    }

    return nil;
    
    } else {
      return fmt.Errorf("Unsure what happened");
    } 
  }();

  err2 := func() error {
    if res.Has_computer {
      res.Has_computer = false;
      res.Issued_computer = "";
      comp.Signed_out_to.Mdoc = 0;

      sqlStatement, prepErr := db.Prepare("UPDATE residents SET has_computer = 0, issued_computer = '' WHERE mdoc = ?");
      if prepErr != nil {
        fmt.Println("Error: Prepare has_computer true ", prepErr);
        return prepErr;
      }

      defer sqlStatement.Close();

      _, execErr := sqlStatement.Exec(res.Mdoc);
      if execErr != nil {
        fmt.Println("Error: Exec has_computer true ", execErr);
        return execErr;
      }

      return nil;

    } else if !res.Has_computer {
      res.Has_computer = true;
      res.Issued_computer = comp.Serial;
      comp.Signed_out_to.Mdoc = res.Mdoc;

      sqlStatement, prepErr := db.Prepare("UPDATE residents SET has_computer = 1, issued_computer = ? WHERE mdoc = ?");
      if prepErr != nil {
        fmt.Println("Error: Prepare has_computer false, ", prepErr);
        return prepErr;
      }

      defer sqlStatement.Close();

      _, execErr := sqlStatement.Exec(comp.Serial, res.Mdoc);
      if execErr != nil {
        fmt.Println("Error: Exec has_computer false, ", execErr);
        return execErr;
      }

      return nil;

    } else {
      fmt.Println("Also not sure what happened")
      return nil;
    } 
  }();   
  
  if err == nil && err2 == nil {
    return nil;
  } else if err != nil {
    return err
  } else {
    return err2;
  }
}

func updateCurrentSignOuts(currentSignOuts []string, data string) []string {
  for i, v := range CurrentSignOuts {
    if v[:31] == data[:31] {
      CurrentSignOuts = append(CurrentSignOuts[:i], CurrentSignOuts[i+1:]...);
      return CurrentSignOuts;
    }
  }
  CurrentSignOuts = append(CurrentSignOuts, data);
  return CurrentSignOuts;
}

func checkComputerStatus(db *structs.Database, serial string, mdoc int) (pass int) {

  check1 := func() int {
    signedout, soErr := db.IsSignedout(serial);
    // if no errors 
    if soErr == nil {
      switch signedout {
      case 0:
        return 1 /* true */;
      case 1:
        signed_out_to, signed_out_to_Err := db.IsSignedoutTo(serial);
        fmt.Println(signed_out_to.Name_of, signed_out_to.Mdoc)
        if signed_out_to_Err == nil {
          if signed_out_to.Mdoc != mdoc {
            return 0 /* false */;
          } else {
            return 1 /* true */;
          }
        } else {
          fmt.Println("Error: IsSignedoutTo failed");
          return -1 /* error */
        }
      default:
        fmt.Println("Error: Default signedout is -1");
        return -1 /* error */
      }
    } else {
      fmt.Println("Error: IsSignedout failed");
      return -1 /* error */
    }
  }();

  check2 := func() int {
    has_computer, hasErr := db.HasComputer(mdoc);
    // if no errors 
    if hasErr == nil {
      switch has_computer {
      case 0:
        return 1 /* true */;
      case 1:
        has_computer_number, has_computer_Err := db.HasComputerNumber(mdoc);
        if has_computer_Err == nil {
          if has_computer_number.Serial != serial {
            if has_computer_number.Serial == "" {
              return 1;
            }
            fmt.Println("==", has_computer_number.Serial)
            return 0 /* false */;
          } else {
            return 1 /* true */;
          }
        } else {
          fmt.Println("Error: IsSignedoutTo failed", has_computer_Err);
          return -1 /* error */
        }
      default:
        fmt.Println("Error: Default signedout is -1");
        return -1 /* error */
      }
    } else {
      fmt.Println("Error: IsSignedout failed");
      return -1 /* error */
    }
  }();

  if check1 == 1 && check2 == 1 {
    return 1;
  } else if check1 == -1 || check2 == -1 {
    fmt.Println(check1, check2);
    return -1;
  } else {
    return 0;
  }

}

// continuing
