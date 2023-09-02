package main

import (
  "Scanner/server"
  "Scanner/sqlite"
  "Scanner/utils"
  "fmt"
  "sync"
)

func main() {
  db, tableErr := sqlite.CreateTables();
  if tableErr != nil {
    fmt.Println("Failed to Create tables.")
    return;
  }
  fmt.Println("Read files")
  utils.ReadResidentsIntoDb(db);

  fmt.Println("Read files 2")
  utils.ReadComputersIntoDb(db);

  fmt.Println("Read files 3")
  var wg sync.WaitGroup;

  wg.Add(1);
  go func() {
    fmt.Println("Server")
    _ = server.Serve("1234");
    wg.Done(); 
  }();

  wg.Wait();

  fmt.Println("ScanScanning")
  utils.ProcessScan(db);

}
