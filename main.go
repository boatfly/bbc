package main

import (
	cli2 "bbc/core/cli"
	"flag"
	"fmt"
)

func main()  {
	cli:=cli2.CLI{}
	cli.RUN()
}

var species = flag.String("species", "go", "the usage of flag")
var num = flag.Int("ins", 1, "ins nums")

func main_202001101() {
	//在flag各种参数生效之前，需要进行解析
	flag.Parse()
	fmt.Println("a string of flag:", *species)
	fmt.Println("ins nums:", *num)
}

func main_20201011() {
	//block := types.NewBlock(1, nil, []byte("first blc block!"))
	//fmt.Printf("the first block in blc is:%v\n", block)

	/// bc := types.CreateBlockChainWithGenesisBlock()
	//fmt.Printf("the first block in blc is:%v\n", bc.Blocks[0])
	//
	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height,bc.Blocks[len(bc.Blocks)-1].Hash,[]byte("Tome send 12 bbc to Alice"))
	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height,bc.Blocks[len(bc.Blocks)-1].Hash,[]byte("Alice send 10 bbc to Plateau"))
	//
	//for _,block:=range bc.Blocks{
	//	fmt.Printf("block info is:%v\n",block)
	//}
	//
	//for _,block:=range bc.Blocks{
	//	fmt.Printf("prehash:%x,hash:%x\n",block.PreBlockHash,block.Hash)
	//}

	///bc.AddBlock([]byte("Tome send 12 bbc to Alice"))
	///bc.AddBlock([]byte("Alice send 10 bbc to Plateau"))

	///bc.PrintChain()

}
