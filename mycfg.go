package main

import (
  "encoding/json"
  "fmt"
  "os"
  "io/ioutil"
)

type Prog struct {
  Functions []Function
}

type Function struct {
  Instrs  []Instruction
  Name    string
}

type Instruction struct {
  Dest    string
  Op      string
  Type    string
  Value   int
  Labels  []string
  Label   string
  Args    []string
}

type Block struct {
  Instrs  []Instruction
}

var Terminators = map[string]bool {
  "jmp": true,
  "ret": true,
  "br": true,
}

func main() {
  var prog Prog
  text, _ := ioutil.ReadFile(os.Args[1])
  json.Unmarshal(text, &prog)

  var blocks []Block
  for _,f := range prog.Functions {
    blocks = create_blocks(f.Instrs)
    block_map, names := create_block_map(blocks)
    cfg := create_cfg(block_map, names)

    graphString := fmt.Sprintf("digraph %s {\n", f.Name)
    for _, name := range names {
      graphString += fmt.Sprintf("  %s;\n", name)
    }
    for _, name := range names {
      for _, succ := range cfg[name] {
        graphString += fmt.Sprintf("  %s -> %s;\n", name, succ)
      }
    }
    graphString += "}"
    fmt.Println(graphString)
  }


}

func create_cfg(block_map map[string]Block, names []string) map[string][]string{

  res := make(map[string][]string)

  for name, block := range(block_map) {
    res[name] = make([]string,0)
    last := block.Instrs[len(block.Instrs)-1]
    if last.Op == "jmp" || last.Op == "br" {
      for _,l := range last.Labels {
        res[name] = append(res[name],l)
      }
    } else if last.Op == "ret" {
      ;
    } else {
      i := 0
      for names[i] != name {
        i += 1
      }
      if i < len(names)-1 {
        res[name] = append(res[name], names[i+1])
      }
    }
  }
  return res
}

func create_blocks(instrs []Instruction) []Block {
  res := make([]Block, 0)
  cur_block := Block{make([]Instruction,0)}

  for _,instr := range instrs {
    // An actual instruction.
    if len(instr.Op) > 0 {
      cur_block.Instrs = append(cur_block.Instrs, instr)

      // Check for terminator.
      if Terminators[instr.Op] == true {
        res = append(res, cur_block)
        cur_block = Block{make([]Instruction,0)}
      }
    // Label.
    } else {
      if len(cur_block.Instrs) > 0 {
        res = append(res, cur_block)
        cur_block = Block{make([]Instruction,0)}
      }
      // Append label to start of new basic block.
      cur_block.Instrs = append(cur_block.Instrs, instr)
    }
  }

  if len(cur_block.Instrs) > 0 {
    res = append(res, cur_block)
  }
  return res
}

func create_block_map(blocks []Block) (map[string]Block, []string) {

  names := make([]string, 0)
  res := make(map[string]Block)
  var name string
  id := 0

  for _,block := range blocks {
    if len(block.Instrs[0].Label) > 0 {
      name = block.Instrs[0].Label
      res[name] = Block{block.Instrs[1:]}
    } else {
      name = fmt.Sprintf("b%d", id)
      id += 1
      res[name] = block
    }
    names = append(names,name)
  }
  return res, names
}






