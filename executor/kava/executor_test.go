// +build integration

package kava

import (
	"math/big"
	"testing"
	"time"

	sdk "github.com/kava-labs/cosmos-sdk/types"
	"github.com/kava-labs/go-sdk/client"
	app "github.com/kava-labs/go-sdk/kava"
	"github.com/stretchr/testify/require"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

func TestSendAmount(t *testing.T) {
	// Note: This test requires kvd to be running locally, with a funded deputy account.

	kavaConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixes(kavaConfig)
	kavaConfig.Seal()

	deputyAddr, err := sdk.AccAddressFromBech32("kava1sl8glhaa9f9tep0d9h8gdcfmwcatghtdrfcd2x")
	require.NoError(t, err)
	coldAddr, err := sdk.AccAddressFromBech32("kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c")
	require.NoError(t, err)
	config := util.KavaConfig{
		KeyType:                    "mnemonic",
		Mnemonic:                   "slab twist stumble inmate predict parent repair crystal celery swarm memory loan rabbit blanket shell talk attend charge inside denial harbor music board steak",
		RpcAddr:                    "tcp://localhost:26657",
		Symbol:                     "bnb",
		DeputyAddr:                 deputyAddr,
		ColdWalletAddr:             coldAddr,
		FetchInterval:              2,
		TokenBalanceAlertThreshold: 10000,
		KavaBalanceAlertThreshold:  10000,
	}
	exe := NewExecutor(client.LocalNetwork, &config)

	amountToSend := big.NewInt(100_000_000)
	previousBalance, err := exe.GetBalance(config.ColdWalletAddr.String())

	txHash, err := exe.SendAmount(config.ColdWalletAddr.String(), amountToSend)
	require.NoError(t, err)

	waitUntilTxInBlock(exe, txHash, 10*time.Second)

	// check coins have moved
	balance, err := exe.GetBalance(config.ColdWalletAddr.String())
	require.NoError(t, err)
	require.Equal(t, balance, previousBalance.Add(previousBalance, amountToSend))
}

func waitUntilTxInBlock(executor common.Executor, txHash string, timeout time.Duration) {
	var status store.TxStatus
	endTime := time.Now().Add(timeout)

	for !(status == store.TxSentStatusSuccess || time.Now().After(endTime)) {
		time.Sleep(500 * time.Millisecond)
		status = executor.GetSentTxStatus(txHash)
	}
}
