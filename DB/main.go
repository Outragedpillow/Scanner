package main

import (
  "Scanner/server"
  "Scanner/sqlite"
  "Scanner/utils"
  "fmt"
)

func main() {
  db, tableErr := sqlite.CreateTables();
  if tableErr != nil {
    fmt.Println("Failed to Create tables.")
    return;
  }
  utils.ReadFromResidents(db);

  if err := server.Serve("1234"); err != nil {
    return;
  }
}
