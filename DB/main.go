package main

import (
  "Scanner/server"
  "Scanner/sqlite"
  "Scanner/utils"
  "fmt" 
)

func main() {
  if deleteErr := utils.DeleteStorageDb(); deleteErr != nil {
    fmt.Println("Error: Deleting Storage")
    return;
  }
  rmSignedoutErr := utils.DeleteSignedout();
  if rmSignedoutErr != nil {
    fmt.Println("Error: Remove signedout.txt ", rmSignedoutErr);
  }

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
