package common

//数据库名称
const DbName = "boat.bbc.db_%s" //根据节点号命名

//表名称
const BlockTableName = "bbc"

//bolt.db中最新区块哈希的key
const LastestHASHKEY = "lastesthash"

//存放UTXO数据
const UTXO_TABLE = "utxos"

//ENV 节点
const NODE_ID = "NODE_ID"

//协议
const RPOTOCOL = "tcp"

//网络服务常量
const COMMAND_LENGTH = 12

//命令分类
const (
	//验证当前节点末端区块是否是最新区块
	CMD_VERSION = "version"
	//从最长链上去获取区块
	CMD_GETBLOCKS = "getblocks"
	//向其他节点展示当前节点有哪些区块
	CMD_INV = "inv"
	//请求指定区块
	CMD_GETDATA = "getdata"
	//接收到新区块，之后进行处理
	CMD_BLOCK = "block"
)
