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
		blocks = form_blocks(f.Instrs)
	}
	for _,block := range blocks {
		fmt.Println(block, string('\n'))
	}
	block_map := create_block_map(blocks)
	fmt.Println(block_map)
}

func form_blocks(instrs []Instruction) []Block {
	res := make([]Block, 0)
	cur_block := Block{make([]Instruction,0)}

	for _,instr := range instrs {
		// An actual instruction.
		if len(instr.Op) > 0 {
			cur_block.instrs = append(cur_block.instrs, instr)

			// Check for terminator.
			if Terminators[instr.Op] == true {
				res = append(res, cur_block)
				cur_block = Block{make([]Instruction,0)}
			}
		// Label.
		} else {
			res = append(res, cur_block)
			cur_block = Block{make([]Instruction,0)}
			// Append label to start of new basic block.
			cur_block.instrs = append(cur_block.instrs, instr)
		}
	}

	if len(cur_block.instrs) > 0 {
		res = append(res, cur_block)
	}
	return res
}

func create_block_map(blocks []Block) map[string]Block {

	res := make(map[string]Block)
	id := 0

	for _,block := range blocks {
		if len(block.instrs[0].Label) > 0 {
			name := block.instrs[0].Label
			res[name] = Block{block.instrs[1:]}
		} else {
			name := fmt.Sprintf("b%d", id)
			id += 1
			res[name] = block
		}
	}
	return res
}






