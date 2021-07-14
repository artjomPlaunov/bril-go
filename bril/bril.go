package bril

import (
  "bytes"
  "encoding/json"
  "strconv"
)

type value string

func (v *value) UnmarshalJSON(data []byte) error {
  if string(data) == "true" || string(data) == "false" {
    *v = value(data)
  } else {
    var tmp int
    if err := json.Unmarshal(data, &tmp); err != nil {
      return err
    }
    *v = value(strconv.Itoa(tmp))
  }

  return nil
}


func (v * value) MarshalJSON() ([]byte, error) {
  buffer := bytes.NewBufferString("")
  buffer.WriteString(string(*v))
  return buffer.Bytes(), nil
}

type Instruction struct {
  Dest    string `json:"dest,omitempty"`
  Op      string `json:"op,omitempty"`
  Type    string `json:"type,omitempty"`
  Value   value `json:"value,omitempty"`
  Labels  []string `json:"labels,omitempty"`
  Label   string `json:"label,omitempty"`
  Args    []string `json:"args,omitempty"`
}

type Arg struct {
  Name  string`json:"name,omitempty"`
  Type  string`json:"type,omitempty"`
}

type Function struct {
  Instrs  []Instruction`json:"instrs"`
  Name    string`json:"name"`
  Args    []Arg`json:"args,omitempty"`
}

type Program struct {
  Functions []Function`json:"functions"`
}

