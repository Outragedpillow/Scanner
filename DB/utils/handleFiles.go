package utils

import ( 
    "bufio"
    "fmt"
    //"log"
    "os"
    "strings"
    "strconv"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "errors"
)

const (
  COMPUTERFILE = "./computers.txt"
  RESIDENT_FILE = "./residents.txt"
  HISTORY_FILE = "./history.txt"
  STORAGE_DB_FILE = "./Storage.db"
  SIGNED_OUT_FILE = "./signedout.txt"
)

func DeleteStorageDb() error {
  noneErr := os.Remove(STORAGE_DB_FILE);
  if noneErr != nil {
    if os.IsNotExist(noneErr) {
      return nil;
    }
    return noneErr;
  }

  return nil;
}


func ReadFilesIntoDb(db *sql.DB) {
  ReadResidentsIntoDb(db);
  ReadComputersIntoDb(db);
}

func ReadComputersIntoDb(db *sql.DB) error {
  
  file, openErr := os.Open(COMPUTERFILE);
  if openErr != nil {
    return openErr;
  }

  scanner := bufio.NewScanner(file);
  
  for scanner.Scan() {
    comp := parseComputers(scanner.Text());
    insertErr := insertComputerData(db, comp);
    if insertErr != nil {
      fmt.Println(insertErr);
    }
  }

  return nil;
}

func parseComputers(info string) []string {
  words := strings.Fields(info);
  if len(words) != 2 {
    fmt.Println("Error: Less than 2");
    return nil;
  }

  if len(words[0]) != 20 {
    if len(words[0]) != 23 {
      fmt.Println("Error: Less than 20");
      return nil;
    }
  } 

  index := strings.Index(words[0], "R");
  if index != -1 {
    words[0] = words[0][index:];
  }

  return words;
}

func insertComputerData(db *sql.DB, info []string) error {
  if len(info) == 2 {
    sqlStatement, prepErr := db.Prepare("INSERT INTO computers (serial, tag_number, is_issued) VALUES (?, ?, ?)");
    if prepErr != nil {
      return prepErr;
    }

    defer sqlStatement.Close();

    _, execErr := sqlStatement.Exec(info[0], info[1], 0);
    if execErr != nil {
      return execErr;
    }

    return nil;
  }

  return errors.New("Error: Insert invalid input length.")
}


func ReadResidentsIntoDb(db *sql.DB) {
  file, openErr := os.Open(RESIDENT_FILE);
  if openErr != nil {
    fmt.Println("Error", openErr)
    return;
  }

  defer file.Close();

  scanner := bufio.NewScanner(file);
  for scanner.Scan() {
    res := parseResidents(scanner.Text());
    insertErr := insertResidentData(db, res);
    if insertErr != nil {
      fmt.Println("Error", insertErr)
      return;
    }
  }
}

func parseResidents(line string) []string {
  words := strings.Fields(line);
  if len(words) != 3 {
    return nil;
  }
  words[0] = strings.Trim(words[0], ",");
  
  return words;
}

func insertResidentData(db *sql.DB, info []string) error {
  mdoc, convErr := strconv.Atoi(info[2]);
  if convErr != nil {
    return convErr;
  }

  name_of := info[0] + " " + info[1];
  sqlStatement, prepErr := db.Prepare("INSERT INTO residents (name_of_r, mdoc) values (?, ?)");
  if prepErr != nil {
    return prepErr;
  }


  defer sqlStatement.Close();

  _, execErr := sqlStatement.Exec(name_of, mdoc);
  if execErr != nil {
    return execErr;
  }

  return nil;
}

func WriteComputerLogs(data string, fileName string) {
  if fileName == "history" {
    file, openErr := os.OpenFile(HISTORY_FILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644);
    if openErr != nil {
      fmt.Println("Error: Open file ", openErr);
    }

    defer file.Close();
     
    _, writeErr := file.WriteString(data);
    if writeErr != nil {
      fmt.Println("Error: Write file, ", writeErr)
    }
  } else if fileName == "signedout" {
      file, openErr := os.OpenFile(SIGNED_OUT_FILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644);
      if openErr != nil {
        fmt.Println("Error: Open file ", openErr);
      }

      defer file.Close();
       
      _, writeErr := file.WriteString(data);
      if writeErr != nil {
        fmt.Println("Error: Write file, ", writeErr)
      }
  }

}
