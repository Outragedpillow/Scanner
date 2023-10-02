package server

import (
	"Scanner/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiData struct {
  CurrentSignOuts []string `json:"currentsignouts"`
}

func SendApiData(w http.ResponseWriter, r *http.Request) {
  data := ApiData {
    CurrentSignOuts: utils.CurrentSignOuts,
  };

  jsonData, jsonErr := json.Marshal(data);
  if jsonErr != nil {
    fmt.Println("Error: Marshal json, ", jsonErr);
    http.Error(w, "Error: Marshal json ", http.StatusInternalServerError);
    return;
  }

  w.Header().Set("Content-Type", "application/json");
  w.WriteHeader(http.StatusOK)
  _, writeErr := w.Write(jsonData)
  if writeErr != nil {
    http.Error(w, "Failed to write response", http.StatusInternalServerError)
    return;
  }

}
