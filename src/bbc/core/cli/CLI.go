package cli

import (
	"bbc/common/utils"
	"bbc/core/types"
	"flag"
	"fmt"
	"log"
	"os"
)

//通过命令行对blockchain进行管理
type CLI struct{}

//参数数量的检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}
}

//用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	//创建钱包
	fmt.Printf("\tcreatewallet -- 创建钱包\n")
	//获取钱包地址列表
	fmt.Printf("\tgetaccounts -- 创建钱包地址集合\n")
	//初始化区块 send -from "[\"boat\"]" -to "[\"plateau\"]" -amount "[\"3\"]"
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
	//余额查询
	fmt.Printf("\tgetbalance -address FROM -- 查询指定地址的余额\n")
	fmt.Printf("\t\t余额查询参数说明:\n")
	fmt.Printf("\t\t\t-address FROM: -- 钱包地址\n")

	//utxo table
	fmt.Printf("\tutxo -method METHOD --测试UTXO Table功能中指定的方法\n")
	fmt.Printf("\t\tMETHOD -- 方法名\n")
	fmt.Printf("\t\t\tbalance -- 查找所有UTXO\n")
	fmt.Printf("\t\t\treset -- 重置UTXO TABLE\n")

	fmt.Printf("\tsetid -port PORT --设置节点号\n")
	fmt.Printf("\tstart --节点启动服务\n")
}

//添加区块
func (cli *CLI) addBlock(txs []*types.Transaction, nodeId string) {
	//判断数据库是否存在？
	if !types.DbExist(nodeId) {
		fmt.Println("数据库不存在！")
		os.Exit(1)
	}
	//获取区块链对象实例
	blockchain := types.BlockChainObject(nodeId)
	blockchain.AddBlock(txs)
}

//命令行运行函数
func (cli *CLI) RUN() {
	nodeId := utils.GetEnvNodeId()
	IsValidArgs()
	//新建相关命令
	//创建钱包
	createwalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getaccountsCmd := flag.NewFlagSet("getaccounts", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)
	createBlockChainWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	//utxo 测试命令行
	flagUtxoCmd := flag.NewFlagSet("utxo", flag.ExitOnError)

	flagsetidCmd := flag.NewFlagSet("setid", flag.ExitOnError)
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)

	//数据参数处理
	flagAddBlockArgs := addBlockCmd.String("data", "send 12 bbc to someone", "添加区块数据的内容")
	//创世区块创建时指定矿工地址（接收奖励）
	flagCreateblockchainArgs := createBlockChainWithGenesisBlockCmd.String("address", "plateau", "指定接收奖励的矿工地址")
	//发起交易参数
	flagSendFromArgs := sendCmd.String("from", "", "转账源地址")
	flagSendToArgs := sendCmd.String("to", "", "转账接收地址")
	flagSendAmountArgs := sendCmd.String("amount", "", "转账金额")

	flagGetbalanceArgs := getbalanceCmd.String("address", "", "钱包地址")

	//utxo table
	flagUtxoCmdArgs := flagUtxoCmd.String("method", "", "UTXO table 相关的操作")

	flagsetidCmdArgs := flagsetidCmd.String("port", "", "节点号")

	//判断命令
	switch os.Args[1] {
	case "start":
		if err := startCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse startCmd failed! %v\n", err)
		}
	case "setid":
		if err := flagsetidCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse flagsetidCmd parameters failed! %v\n", err)
		}
	case "getaccounts":
		if err := getaccountsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse getaccountsCmd parameters failed! %v\n", err)
		}
	case "createwallet":
		if err := createwalletCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createwalletCmd parameters failed! %v\n", err)
		}
	case "getbalance":
		if err := getbalanceCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse getbalanceCmd parameters failed! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse sendCmd parameters failed! %v\n", err)
		}
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
	case "utxo":
		if err := flagUtxoCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse flagUtxoCmd parameters failed! %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	//节点启动服务
	if startCmd.Parsed() {
		cli.startNode(nodeId)
	}

	//设置节点号
	if flagsetidCmd.Parsed() {
		cli.SetNodeId(*flagsetidCmdArgs)
	}

	//创建钱包
	if createwalletCmd.Parsed() {
		cli.createwallets(nodeId)
	}
	//获取地址集合
	if getaccountsCmd.Parsed() {
		cli.getAccounts(nodeId)
	}

	//utxo cli
	if flagUtxoCmd.Parsed() {
		if *flagUtxoCmdArgs == "" {
			fmt.Println("UTXO method不能为空")
			PrintUsage()
			os.Exit(1)
		}
		switch *flagUtxoCmdArgs {
		case "balance":
			cli.TestFindUtxoMap()
		case "reset":
			cli.TestResetUTXO(nodeId)
		default:

		}
	}

	//查询余额
	if getbalanceCmd.Parsed() {
		if *flagGetbalanceArgs == "" {
			fmt.Println("钱包地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetbalanceArgs, nodeId)
	}

	//发起转账
	if sendCmd.Parsed() {
		if *flagSendFromArgs == "" {
			fmt.Println("转账源地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArgs == "" {
			fmt.Println("转账接收地址")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArgs == "" {
			fmt.Println("转账金额不能为空")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFrom[%s]\n", utils.JSON2Slice(*flagSendFromArgs))
		fmt.Printf("\tTo[%s]\n", utils.JSON2Slice(*flagSendToArgs))
		fmt.Printf("\tAmount[%s]\n", utils.JSON2Slice(*flagSendAmountArgs))
		cli.send(utils.JSON2Slice(*flagSendFromArgs), utils.JSON2Slice(*flagSendToArgs), utils.JSON2Slice(*flagSendAmountArgs), nodeId)
	}

	if printBlockChainCmd.Parsed() {
		cli.printBlockChain(nodeId)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockArgs == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*types.Transaction{}, nodeId)
	}

	if createBlockChainWithGenesisBlockCmd.Parsed() {
		if *flagCreateblockchainArgs == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateblockchainArgs, nodeId)
	}

}
