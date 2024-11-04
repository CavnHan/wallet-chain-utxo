package bitcoin

import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/btcsuite/btcd/rpcclient"
)

type BtcClient struct {
	*rpcclient.Client
	compressed bool
}

func NewBtcClient(RpcUrl, RpcUser, RpcPass string) (*BtcClient, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         RpcUrl,
		User:         RpcUser,
		Pass:         RpcPass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}, nil)
	if err != nil {
		log.Error("new bitcoin rpc client fail", "err:", err)
		return nil, err
	}
	return &BtcClient{
		Client:     client,
		compressed: true,
	}, nil
}

func (bc *BtcClient) GetLatestBlock() int64 {
	count, err := bc.Client.GetBlockCount()
	if err != nil {
		return 0
	}
	return count

}
