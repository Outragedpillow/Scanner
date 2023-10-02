package server

import (
  "net/http"
)

func Serve(port string) error {
  // id := "12";
  // http.HandleFunc("/api/sign-out/"+id, api.FindResident);
  
  http.HandleFunc("/api/data", SendApiData);

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    return err;
  }

  return nil;
}
