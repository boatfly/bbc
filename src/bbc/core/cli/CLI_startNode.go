package cli

import "bbc/core/types"

//节点启动服务
func (cli *CLI) startNode(nodeId string) {
	types.StartServer(nodeId)
}
