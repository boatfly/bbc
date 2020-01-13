package cli

import (
	"bbc/core/types"
	"flag"
	"fmt"
	"log"
	"os"
)

//通过命令行对blockchain进行管理
type CLI struct {}

//用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	//初始化区块
	fmt.Printf("\tcreateblockchain -address ADDRESS -- 创建区块链\n")
	//添加区块
	fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	//打印完整的区块信息
	fmt.Printf("\tprintblockchain -- 打印完整区块链信息\n")
	//通过命令行转账
	fmt.Printf("\tsned -from FROM -to TO -amount AMOUNT --发起转账\n")
	fmt.Printf("\t\t转账参数说明:\n")
	fmt.Printf("\t\t\t-from FROM: -- 转账源地址\n")
	fmt.Printf("\t\t\t-to TO: -- 转账接收地址\n")
	fmt.Printf("\t\t\t-amount AMOUNT: -- 转账金额\n")
}

//参数数量的检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}
}

//发起交易
func (cli *CLI) send()  {
	
}

//初始化区块
func (cli *CLI) createBlockChain(address string) {
	types.CreateBlockChainWithGenesisBlock(address)
}

//添加区块
func (cli *CLI) addBlock(txs []*types.Transaction) {
	//判断数据库是否存在？
	if !types.DbExist() {
		fmt.Println("数据库不存在！")
		os.Exit(1)
	}
	//获取区块链对象实例
	blockchain := types.BlockChainObject()
	blockchain.AddBlock(txs)
}

//打印完整的区块信息
func (cli *CLI) printBlockChain() {
	//判断数据库是否存在？
	if !types.DbExist() {
		fmt.Println("数据库不存在！")
		os.Exit(1)
	}
	//获取区块链对象实例
	blockchain := types.BlockChainObject()
	blockchain.PrintChain()
}

//命令行运行函数
func (cli *CLI) RUN() {
	IsValidArgs()
	//新建相关命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)
	createBlockChainWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	//sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	//数据参数处理
	flagAddBlockArgs := addBlockCmd.String("data", "send 12 bbc to someone", "添加区块数据的内容")
	//创世区块创建时指定矿工地址（接收奖励）
	flagCreateblockchainArgs := createBlockChainWithGenesisBlockCmd.String("address", "plateau", "指定接收奖励的矿工地址")

	//判断命令
	switch os.Args[1] {
	case "send":

	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse addBlockCmd parameters failed! %v\n", err)
		}
	case "printblockchain":
		if err := printBlockChainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse printBlockChainCmd parameters failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBlockChainWithGenesisBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBlockChainWithGenesisBlockCmd parameters failed! %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	if printBlockChainCmd.Parsed() {
		cli.printBlockChain()
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockArgs == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*types.Transaction{})
	}

	if createBlockChainWithGenesisBlockCmd.Parsed() {
		if *flagCreateblockchainArgs == ""{
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateblockchainArgs)
	}

}
