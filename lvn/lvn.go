package main

import (
  "bril"
  "cfg"
  "encoding/json"
  "fmt"
  "os"
  "io/ioutil"
  "strconv"
)

// If Op is "const", then First holds the constant value.
type Lvn_value struct {
  Op      string
  First   int
  Second  int
}

type Lvn_row struct {
  Num     int
  Val_num string
}

func main() {
  var prog bril.Program
  text, _ := ioutil.ReadFile(os.Args[1])
  json.Unmarshal(text, &prog)
  var blocks []cfg.Block
  blocks = cfg.Create_blocks(prog.Functions[0].Instrs)

  for i,b := range blocks {
    blocks[i] = lvn(b)
  }

  res,_ := json.Marshal(prog)
  fmt.Println(string(res))
}

func lvn(block cfg.Block) cfg.Block {
  // New var names will have the form "lvn.<id>", i.e., lvn.0, lvn.1, etc. 
  //next_var_id := 0

  // Lvn Table mapping canonical values to their value number and variable name.
  //lvn_table := make(map[Value]Lvn_row)

  // Environment mapping variable names to their current value numbers.
  env := make(map[string]int)

  for _, instr := range block.Instrs {
    lvn_val := construct_lvn_val(instr, env)
    // Handle print instructions
    if instr.Op == "print" {

    } else {
      
    }
  }
  return block
}

func construct_lvn_val(instr bril.Instruction, env map[string]int) Lvn_value {
  // Const Instructions.
  if instr.Op == "const" {
    tmp, _ := strconv.Atoi(string(instr.Value))
    // If Op is const, then the First struct field holds the constant value.
    return Lvn_value{instr.Op,tmp,0}
  // Add, Mul, Sub, and Div Instructions.
  } else if ( instr.Op == "add" ||
              instr.Op == "mul" ||
              instr.Op == "sub" ||
              instr.Op == "div") {
    return Lvn_value{instr.Op, env[instr.Args[0], env[instr.Args[1]}
  }
  return Lvn_value{}
}







