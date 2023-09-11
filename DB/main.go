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

  utils.ReadFilesIntoDb(db);

  go func() {
    serverErr := server.Serve("1234");
    if serverErr != nil {
      return;
    }
  }()

  go func() {
    utils.ProcessScan(db);
  }();

  select{};
}
