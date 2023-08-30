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
)

func ReadResidentsIntoDb(db *sql.DB) {
  file, openErr := os.Open("./utils/residents.txt");
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
  words := strings.Split(line, " ");
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
  sqlStatement, prepErr := db.Prepare("INSERT INTO residents (name_of, mdoc) values (?, ?)");
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


