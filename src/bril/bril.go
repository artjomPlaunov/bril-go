package bril

type Instruction struct {
  Dest    string
  Op      string
  Type    string
  Value   int
  Labels  []string
  Label   string
  Args    []string
}

type Function struct {
  Instrs  []Instruction
  Name    string
}

type Program struct {
  Functions []Function
}


