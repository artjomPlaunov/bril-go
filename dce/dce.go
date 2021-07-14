package main

import (
  "bril"
  "cfg"
  "encoding/json"
  "fmt"
  "os"
  "io/ioutil"
)

func main() {
  var prog bril.Program
  text, _ := ioutil.ReadFile(os.Args[1])
  json.Unmarshal(text, &prog)
  var blocks []cfg.Block
  blocks = cfg.Create_blocks(prog.Functions[0].Instrs)


  // Global DCE
  blocks = dce_global(blocks)

  // Local DCE
  for i,b := range blocks {
    blocks[i] = dce_local(b)
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


func dce_global(blocks []cfg.Block) []cfg.Block {
  changed, blocks := dce_global_aux(blocks)

  for changed {
    changed, blocks = dce_global_aux(blocks)
  }
  return blocks
}

func dce_global_aux(blocks []cfg.Block) (bool, []cfg.Block) {
  used := make(map[string]bool)
  for _,b := range blocks {
    for _,instr := range b.Instrs {
      if len(instr.Args) > 0 {
        for _,arg := range instr.Args {
          used[arg] = true
        }
      }
    }
  }

  for i,b := range blocks {
    for j,instr := range b.Instrs {
      if len(instr.Dest) > 0 {
        _,ok := used[instr.Dest]
        if !ok {
          blocks[i].Instrs = append(blocks[i].Instrs[:j],
                                     blocks[i].Instrs[j+1:]...)
          return true, blocks
        }
      }
    }
  }

  return false, blocks
}


func dce_local(b cfg.Block) cfg.Block {
  // Iterate dce algorithim (in dce_aux) to convergence.
  changed, b := dce_local_aux(b)

  for changed {
    changed, b = dce_local_aux(b)
  }
  return b
}

func dce_local_aux(b cfg.Block) (bool, cfg.Block) {
  last_def := make(map[string]int)

  for i, instr := range b.Instrs {
    // Check for uses
    if len(instr.Args) > 0 {
      for _, arg := range instr.Args {
        _, ok := last_def[arg]
        if ok {
          delete(last_def, arg)
        }
      }
    }
    if len(instr.Dest) > 0 {
      idx, ok := last_def[instr.Dest]
      if ok {
        res := append(b.Instrs[:idx], b.Instrs[idx+1:]...)
        b.Instrs = res
        return true, b
      }
      last_def[instr.Dest] = i
    }
  }
  return false, b
}

func print_blocks(blocks []cfg.Block, s string) {
  fmt.Println(s)
  for _,b := range blocks {
    for _,i := range b.Instrs {
      fmt.Println(i)
    }
    fmt.Println("\n")
  }
}












