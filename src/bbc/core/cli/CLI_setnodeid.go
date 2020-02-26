package cli

import (
	"bbc/common"
	"fmt"
	"log"
	"os"
)

//设置端口号
func (cli *CLI) SetNodeId(nodeId string) {
	if nodeId == "" {
		fmt.Println("please set the node id...")
		os.Exit(1)
	}
	err := os.Setenv(common.NODE_ID, nodeId)
	if nil != err {
		log.Fatalf("set env failed!%v\n", err)
	}
	fmt.Println("node_id:"+os.Getenv(common.NODE_ID))
}
