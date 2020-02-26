package cli

import "bbc/core/types"

//重置UTXO TABLE
func (cli *CLI) TestResetUTXO(nodeId string) {
	bc:=types.BlockChainObject(nodeId)
	defer bc.DB.Close()
	utxoSet:=types.UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}

//查找
func (cli *CLI) TestFindUtxoMap()  {


}
