package types

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//钱包集合持久化文件
const WALLETS_FILE = "/Users/song/go/boatWallets_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet //key:address value:Wallet
}

//初始化钱包集合
func NewWallets(nodeId string) *Wallets {
	walletfile:=fmt.Sprintf(WALLETS_FILE,nodeId)
	//从钱包文件中读取钱包
	if _, err := os.Stat(walletfile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets
	}
	fileConent, err := ioutil.ReadFile(walletfile)
	if nil != err {
		log.Panicf("read the wallets file failed!%v\n", err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileConent))
	err = decoder.Decode(&wallets)
	if nil != err {
		log.Panicf("decode the wallets failed!%v\n", err)
	}
	return &wallets
}

//添加新的钱包到集合中
func (wallets *Wallets) CreateWallet(nodeId string) *Wallet {
	wallet := NewWallet()                                 //创建钱包
	wallets.Wallets[string(wallet.GetAddress())] = wallet //加入钱包集合
	//持久化钱包
	wallets.SaveWallets(nodeId)
	return wallet
}

//持久化钱包信息(存储到文件中)
func (w *Wallets) SaveWallets(nodeId string) {
	var content bytes.Buffer //钱包内容
	//注册p256椭圆，注册之后，可以直接在内部对curve的接口进行编码
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if nil != err {
		log.Panicf("wallets encode failed！%v\n", err)
	}
	walletfile:=fmt.Sprintf(WALLETS_FILE,nodeId)
	err = ioutil.WriteFile(walletfile, content.Bytes(), 0644) //这种方式是如果file不存在，创建；如果文件存在，先把原来的file清空，在写入，效率低@TODO
	if nil != err {
		log.Panicf("write the content of wallet into file[%s] failed!%v\n", WALLETS_FILE, err)
	}

}
