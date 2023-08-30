package utils

import ( 
    "bufio"
    "fmt"
    //"log"
    "os"
    "strings"
    // "strconv"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "errors"
)

func ReadComputersIntoDb(db *sql.DB) error {
  file, openErr := os.Open("./utils/computers.txt");
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
  words := strings.Split(info, " ");
  if len(words) != 2 {
    return nil;
  }

  if len(words[0]) != 20 {
    return nil;
  } 

  index := strings.Index(words[0], "R");
  if index != -1 && index < len(words[0])-1 {
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

  return errors.New("Invalid input length.")
}
