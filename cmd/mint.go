package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"erc20-go-demo/internal/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	cfg := loadCfg()
	client := dialClient(cfg.RPCURL)
	defer client.Close()

	auth, _ := signer(cfg)

	contractAddr := common.HexToAddress("<PUT-YOUR-CONTRACT-ADDRESS>")
	erc20, err := token.NewMyToken(contractAddr, client)
	if err != nil {
		log.Fatalf("bind err: %v", err)
	}

	to := common.HexToAddress(cfg.ToAddress) // 也可填自己
	amount := new(big.Int)
	// 示例：增发 10 * 10^18
	amount.SetString("10000000000000000000", 10)

	tx, err := erc20.Mint(auth, to, amount)
	if err != nil {
		log.Fatalf("mint failed: %v", err)
	}
	fmt.Printf("🪙 Mint submitted: %s\n", tx.Hash().Hex())

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("wait mined: %v", err)
	}
	fmt.Println("✅ Mint mined.")
}
