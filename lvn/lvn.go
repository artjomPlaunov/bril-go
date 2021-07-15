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
  Name    string
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
  // Local Value Numbers
  next_lvn := 1

  // New var names will have the form "lvn.<id>", i.e., lvn.0, lvn.1, etc. 
  next_var_id := 0

  // Lvn Table mapping canonical values to their value number and variable name.
  lvn_table := make(map[Value]Lvn_row)

  // Environment mapping variable names to their current value numbers.
  env := make(map[string]int)

  for i, instr := range block.Instrs {
    lvn_val := construct_lvn_val(instr, env)
    // Handle print instructions.
    if instr.Op == "print" {
      // <><><><><><><><><>
    // Handle arithmetic instructions.
    } else if is_arith_op(instr.Op) {
      // Value has been computed before; reuse it.
      if lvn_row, ok := lvn_table[lvn_val]; ok {
        env[instr.Dest] = lvn_row.Num
        // Replace instruction with id instruction.
        block.Instrs[i] = construct_id_instr(lvn_row.Name, instr.Dest)
      } else {


      }
    } else if instr.Op == "const" {
      if lvn_row, ok := lvn_table[lvn_val]; ok {
        env[instr.Dest] = lvn_row.Num
        block.Instrs[i] = construct_id_instr(lvn_row.Name, instr.Dest)
      } else {
        var dest string
        if is_overwritten(block.Instrs, i, instr.Dest) {
          dest = "lvn." + strconv.Itoa(next_var_id)
          next_var_id += 1
        } else {
          dest = instr.Dest
        }
        lvn_table[lvn_val] = Lvn_row{next_lvn, dest}
        for j,arg := range instr.Args {
          
        }
      }
    }
  }
  return block
}

// Helper constructor functions

func construct_id_instr(name, dest string) bril.Instruction {
  return  bril.instruction{
            Op:   "id",
            Args: []string{name},
            Type: "int",
            Dest: dest }
}

func construct_lvn_val(instr bril.Instruction, env map[string]int) Lvn_value {
  // Const Instructions.
  if instr.Op == "const" {
    tmp, _ := strconv.Atoi(string(instr.Value))
    // If Op is const, then the First struct field holds the constant value.
    return Lvn_value{instr.Op,tmp,0}
  // Add, Mul, Sub, and Div Instructions.
  } else if is_arith_op(instr.Op) {
    return Lvn_value{instr.Op, env[instr.Args[0], env[instr.Args[1]}
  }
  }
  return Lvn_value{}
}

// Helper Conditional checks

func is_arith_op(op string) bool {
  if op == "add" || op == "mul" || op == "sub" || op = "div" {
    return true
  } else {
    return false
  }
}
func is_overwritten(instrs []bril.Instruction, idx int, dest string) {
  for i := idx+1; i < len(instrs); i++ {
    if instrs[i].Dest == dest {
      return true
    }
  }
  return false
}





