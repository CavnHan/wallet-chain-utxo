package bitcoin

import (
	"time"

	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
	"github.com/ethereum/go-ethereum/log"
)

type BitcoinData struct {
	BitcoinDataCli *oklink.ChainExplorerAdaptor
}

func NewBitcoinDataClient(baseUrl, apiKey string) (*BitcoinData, error) {
	btcCli, err := oklink.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Second*15)
	if err != nil {
		log.Error("New bitcoin client fail", "err:", err)
		return nil, err
	}
	return &BitcoinData{
		BitcoinDataCli: btcCli,
	}, err
}

func (bd *BitcoinData) GetFee() error {
	return nil
}

func (bd *BitcoinData) GetTxUtxoList() error {
	return nil
}

func (bd *BitcoinData) GetTxByAddress() error {
	return nil
}
