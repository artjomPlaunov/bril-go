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
  Value   Lvn_value
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
  instrs := make([]bril.Instruction, 0)
  for _,b := range blocks {
    for _,i := range b.Instrs {
      instrs = append(instrs, i)
    }
  }
  prog.Functions[0].Instrs = instrs

  res,_ := json.Marshal(prog)
  fmt.Println(string(res))
}

// lvn algorithm
func lvn(block cfg.Block) cfg.Block {
  // Local Value Numbers
  next_lvn := 1

  // New var names will have the form "lvn.<id>", i.e., lvn.0, lvn.1, etc. 
  next_var_id := 0

  // Lvn Table mapping canonical values to their value number and variable name.
  lvn_table := make([]Lvn_row, 0)

  // Environment mapping variable names to their current value numbers.
  env := make(map[string]int)

  for i, instr := range block.Instrs {
    lvn_val,ok := construct_lvn_val(instr, env)
    // Short-circuit and return on blocks that try to use non-local values. 
    if !ok {
      return block
    }

    var lvn int
    var dest string
    if instr.Op == "id" {
      dest = check_overwrite(block, i, instr.Dest, &next_var_id)
      instr.Dest = dest
      env[dest] = env[instr.Args[0]]
      block.Instrs[i] = construct_id_instr(
                          get_var(lvn_table, env[dest]),
                          dest)
    } else if len(instr.Dest) > 0 {
      if lvn_row, ok := contains_val(lvn_table, lvn_val); ok {
        lvn = lvn_row.Num
        block.Instrs[i] = construct_id_instr(lvn_row.Name, instr.Dest)
      } else {
        dest = check_overwrite(block, i, instr.Dest, &next_var_id)
        lvn_table = append(lvn_table, Lvn_row{next_lvn, lvn_val ,dest})
        lvn = next_lvn
        next_lvn += 1
        for j,arg := range instr.Args {
          instr.Args[j] = get_var(lvn_table, env[arg])
        }
      }
      env[instr.Dest] = lvn
    }
  }/*
  for _,row := range lvn_table {
    fmt.Println(row)
  }
  fmt.Println(env)*/
  return block
}

func get_var(lvn_table []Lvn_row, lvn int) string {
  for _, row := range lvn_table {
    if row.Num == lvn {
      return row.Name
    }
  }
  return ""
}

// Lvn Table functions.
func contains_val(lvn_table []Lvn_row, lvn_val Lvn_value) (Lvn_row, bool) {

  for _,row := range lvn_table {
    if is_eq_lvn_val(row.Value, lvn_val) {
      return row,true
    }
  }
  return Lvn_row{}, false
}

// Helper constructor functions.

func construct_id_instr(name, dest string) bril.Instruction {
  return  bril.Instruction{
            Op:   "id",
            Args: []string{name},
            Type: "int",
            Dest: dest }
}

func construct_lvn_val(instr bril.Instruction,
                       env map[string]int) (Lvn_value,bool) {
  // Const Instructions.
  if instr.Op == "const" {
    tmp, _ := strconv.Atoi(string(instr.Value))
    // If Op is const, then the First struct field holds the constant value.
    return Lvn_value{instr.Op,tmp,0}, true
  // Add, Mul, Sub, and Div Instructions.
  } else if is_arith_op(instr.Op) {
    v1,ok1 := env[instr.Args[0]]
    v2,ok2 := env[instr.Args[1]]
    if (!ok1) || (!ok2) {
      return Lvn_value{}, false
    }
    return Lvn_value{instr.Op, v1, v2}, true
  } else if instr.Op == "id" {
    v,ok := env[instr.Args[0]]
    if (!ok) {
      return Lvn_value{}, false
    }
    return Lvn_value{instr.Op, v, 0}, true
  }
  return Lvn_value{}, true
}

// Helper Conditional checks

func is_arith_op(op string) bool {
  if op == "add" || op == "mul" || op == "sub" || op == "div" {
    return true
  } else {
    return false
  }
}
func is_overwritten(instrs []bril.Instruction, idx int, dest string) bool {
  for i := idx+1; i < len(instrs); i++ {
    if instrs[i].Dest == dest {
      return true
    }
  }
  return false
}

func is_eq_lvn_val(v1, v2 Lvn_value) bool {
  if v1.Op == v2.Op && v1.First == v2.First && v1.Second == v2.Second {
    return true
  } else {
    return false
  }
}

func check_overwrite(block cfg.Block,
                     i int,
                     dest string,
                     next_var_id * int) string {
  if is_overwritten(block.Instrs, i, dest) {
    res := "lvn." + strconv.Itoa(*next_var_id)
    *next_var_id += 1
    block.Instrs[i].Dest = res
    return res
  } else {
    return dest
  }
}
