package bitcoin

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/log"

	"github.com/CavnHan/wallet-chain-utxo/chain"
	"github.com/CavnHan/wallet-chain-utxo/config"
	common "github.com/CavnHan/wallet-chain-utxo/rpc/common"
	"github.com/CavnHan/wallet-chain-utxo/rpc/utxo"
)

const ChainName = "Bitcoin"

type ChainAdaptor struct {
	btcClient     *BtcClient
	btcDataClient *BitcoinData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	BtcClient, err := NewBtcClient(conf.WalletNode.Btc.RpcUrl, conf.WalletNode.Btc.RpcUser, conf.WalletNode.Btc.RpcPass)
	if err != nil {
		log.Error("new bitcoin rpc client fail", "err:", err)
		return nil, err
	}
	btcDataClient, err := NewBitcoinDataClient(conf.WalletNode.Btc.DataApiUrl, conf.WalletNode.Btc.DataApiKey)
	if err != nil {
		log.Error("new bitcoin data client fail", "err:", err)
		return nil, err
	}
	return &ChainAdaptor{
		btcClient:     BtcClient,
		btcDataClient: btcDataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	return &utxo.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	switch req.Format {
	case "p2pkh":
		return nil, nil
	case "p2wpkh":
		return nil, nil
	case "p2sh":
		return nil, nil
	case "p2tr":
		return nil, nil
	default:
		return nil, nil
	}
}

func (c *ChainAdaptor) ValidAddress(req *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	switch req.Format {
	case "p2pkh":
		return nil, nil
	case "p2wpkh":
		return nil, nil
	case "p2sh":
		return nil, nil
	case "p2tr":
		return nil, nil
	default:
		return nil, nil
	}
}

func (c *ChainAdaptor) GetFee(req *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetAccount(req *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetUnspentOutputs(req *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockByNumber(req *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	//通过区块高度获取区块hash
	blockHash, err := c.btcClient.Client.GetBlockHash(req.Height)
	if err != nil {
		log.Error("get block hash fail", "err:", err)
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block hash fail",
		}, err
	}
	//初始化一个json.RawMessage类型的切片
	var params []json.RawMessage
	numBlockJson, err := json.Marshal(blockHash)
	if err != nil {
		log.Error("marshal block hash fail", "err:", err)
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "marshal block hash fail",
		}, err
	}
	params = []json.RawMessage{numBlockJson}
	//调用rpcclient的Call方法
	block, err := c.btcClient.Client.RawRequest("getblock", params)
	if err != nil {
		log.Error("get block fail", "err:", err)
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block fail",
		}, err
	}
	var resultBlock BlockData
	//解析json数据
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("get raw transaction fail", "err:", err)
	}
	//处理交易数据
	for _, txid := range resultBlock.Tx {
		txIdJson, _ := json.Marshal(txid)
		boolJson, _ := json.Marshal(true)
		//将txid和boolJson转换为json.RawMessage类型
		dataJson := []json.RawMessage{txIdJson, boolJson}
		//调用rpcclient的Call方法
		tx, err := c.btcClient.Client.RawRequest("getrawtransaction", dataJson)
		if err != nil {
			fmt.Println("get raw transaction fail", "err", err)
		}
		var rawTx RawTransactionData
		//解析json数据
		err = json.Unmarshal(tx, &rawTx)
		if err != nil {
			log.Error("json unmarshal fail", "err:", err)
			return nil, err
		}
		for _, v := range rawTx.Vin {
			fmt.Println("v.TxId==", v.TxId)
		}
	}
	return &utxo.BlockResponse{}, err

}

func (c *ChainAdaptor) GetBlockByHash(req *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	var params []json.RawMessage
	numBlocksJSON, _ := json.Marshal(req.Hash)
	params = []json.RawMessage{numBlocksJSON}
	block, _ := c.btcClient.Client.RawRequest("getblock", params)
	var resultBlock BlockData
	err := json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("Unmarshal json fail", "err", err)
	}
	for _, txid := range resultBlock.Tx {
		txIdJson, _ := json.Marshal(txid)
		boolJSON, _ := json.Marshal(true)
		dataJSON := []json.RawMessage{txIdJson, boolJSON}
		tx, err := c.btcClient.Client.RawRequest("getrawtransaction", dataJSON)
		if err != nil {
			fmt.Println("get raw transaction fail", "err", err)
		}
		var rawTx RawTransactionData
		err = json.Unmarshal(tx, &rawTx)
		if err != nil {
			log.Error("json unmarshal fail", "err", err)
			return nil, err
		}
		for _, v := range rawTx.Vin {
			fmt.Println("v.TxId==", v.TxId)
		}

	}
	return &utxo.BlockResponse{}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	//将hash字符串转换为hash类型
	hash, err := chainhash.NewHashFromStr(req.Hash)
	if err != nil {
		log.Error("format string to hash fail", "err", err)
	}
	blockHeader, err := c.btcClient.Client.GetBlockHeader(hash)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}
	return &utxo.BlockHeaderResponse{
		Code:       common.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: blockHeader.PrevBlock.String(),
		Number:     string(blockHeader.Version),
		Blockhash:  req.Hash,
		Merkleroot: blockHeader.MerkleRoot.String(),
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	blockNumber := req.Height
	if req.Height == 0 {
		latestBlock, err := c.btcClient.Client.GetBlockCount()
		if err != nil {
			return &utxo.BlockHeaderResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "get latest block fail",
			}, err
		}
		blockNumber = latestBlock
	}
	blockHash, err := c.btcClient.Client.GetBlockHash(blockNumber)
	if err != nil {
		log.Error("get block hash by number fail", "err", err)
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block hash fail",
		}, err
	}
	blockHeader, err := c.btcClient.Client.GetBlockHeader(blockHash)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}
	return &utxo.BlockHeaderResponse{
		Code:       common.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: blockHeader.PrevBlock.String(),
		Number:     strconv.FormatInt(blockNumber, 10),
		Blockhash:  blockHash.String(),
		Merkleroot: blockHeader.MerkleRoot.String(),
	}, nil
}

func (c *ChainAdaptor) SendTx(req *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByAddress(req *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByHash(req *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) BuildSignedTransaction(req *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) DecodeTransaction(req *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}
