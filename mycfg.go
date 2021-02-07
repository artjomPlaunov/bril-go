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

  for _,f := range prog.Functions {
    blocks = cfg.Create_blocks(f.Instrs)
    block_map, names := cfg.Create_block_map(blocks)
    mycfg := cfg.Create_cfg(block_map, names)

    graphString := fmt.Sprintf("digraph %s {\n", f.Name)
    for _, name := range names {
      graphString += fmt.Sprintf("  %s;\n", name)
    }
    for _, name := range names {
      for _, succ := range mycfg[name] {
        graphString += fmt.Sprintf("  %s -> %s;\n", name, succ)
      }
    }
    graphString += "}"
    fmt.Println(graphString)
  }
}
