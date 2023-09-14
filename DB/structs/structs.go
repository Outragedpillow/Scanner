package structs

type Admin struct {
  Name_of string `json:"name_of"`
}

type Resident struct {
  Mdoc int `json:"mdoc"`
  Name_of string `json:"name_of"`
}

type Computer struct {
  Serial string `json:"serial"`
  Tag_number int `json:"tag_number"`
  Is_issued bool `json:"is_issued"`
  Signed_out_by Admin `json:"signed_out_by"`
  Signed_out_to Resident `json:"signed_out_to"`
  Time_issued string `json:"time_issued"`
  Time_returned string `json:"time_returned"`
}

func ResidentIsEmpty(s Resident) bool {
  if s.Mdoc == 0 || s.Name_of == "" {
    return true;
  }

  return false;
}

func ComputerIsEmpty(c Computer) bool {
  if c.Serial == "" || c.Tag_number == 0 {
    return true;
  }

  return false;
}
