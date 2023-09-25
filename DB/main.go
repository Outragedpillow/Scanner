package main

import (
	"Scanner/server"
	"Scanner/structs"
	"Scanner/utils"
	"fmt"
)

func main() {
  var db structs.Database;

  checkErr := utils.CheckForCrash();
  if checkErr != nil {
    if deleteErr := utils.DeleteStorageDb(); deleteErr != nil {
        fmt.Println("Error: Deleting Storage")
        return;
      }
      rmSignedoutErr := utils.DeleteSignedout();
      if rmSignedoutErr != nil {
        fmt.Println("Error: Remove signedout.txt ", rmSignedoutErr);
      }

    openDbErr := db.Open("Storage.db");
    if openDbErr != nil {
      panic(openDbErr);
    }

    defer db.Close();

    createErr := db.CreateTables();
    if createErr != nil {
      panic(createErr);
    }
     
    utils.ReadFilesIntoDb(db.Conn);
  }

  openDbErr := db.Open("Storage.db");
    if openDbErr != nil {
      panic(openDbErr);
    }

    defer db.Close();
 
  go func() {
    serverErr := server.Serve("1234");
    if serverErr != nil {
      return;
    }
  }()

  go func() {
    utils.ProcessScan(&db, db.Conn);
  }();

  select{};
}
