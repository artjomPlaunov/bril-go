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
	Instrs	[]Instruction
	Name		string
}

type Instruction struct {
	Dest		string
	Op			string
	Type		string
	Value		int
	Labels	[]string
	Label		string
	Args		[]string
}

type Block struct {
	instrs	[]Instruction
}

func main() {
	var prog Prog
	text, _ := ioutil.ReadFile(os.Args[1])
	json.Unmarshal(text, &prog)

	var blocks []Block
	for _,f := range prog.Functions {
		blocks = form_blocks(f.Instrs)
	}
	fmt.Println(blocks)

}

func form_blocks(instrs []Instruction) []Block {
	res := make([]Block, 0)
	cur_block := Block{make([]Instruction,0)}

	for _,instr := range instrs {
		// An actual instruction.
		if len(instr.Op) > 0 {
			cur_block = append(cur_block, instr)

			// Check for terminator.
			 
		// Label
		} else {

		}
	}
}






