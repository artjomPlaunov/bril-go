package cfg

import (
  "bril"
  "fmt"
)

type Block struct {
  Instrs  []bril.Instruction
}

var Terminators = map[string]bool {
  "jmp": true,
  "ret": true,
  "br": true,
}

func Create_cfg(block_map map[string]Block, names []string) map[string][]string{

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

func Create_blocks(instrs []bril.Instruction) []Block {
  res := make([]Block, 0)
  cur_block := Block{make([]bril.Instruction,0)}

  for _,instr := range instrs {
    // An actual instruction.
    if len(instr.Op) > 0 {
      cur_block.Instrs = append(cur_block.Instrs, instr)

      // Check for terminator.
      if Terminators[instr.Op] == true {
        res = append(res, cur_block)
        cur_block = Block{make([]bril.Instruction,0)}
      }
    // Label.
    } else {
      if len(cur_block.Instrs) > 0 {
        res = append(res, cur_block)
        cur_block = Block{make([]bril.Instruction,0)}
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

func Create_block_map(blocks []Block) (map[string]Block, []string) {

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


