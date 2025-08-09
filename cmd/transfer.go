package main

import (
	"context"
	"fmt"
	"log"

	"erc20-go-demo/internal/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	cfg := loadCfg()
	if cfg.ToAddress == "" || cfg.TransferWei.Sign() == 0 {
		log.Fatalf("need TO_ADDRESS and TRANSFER_AMOUNT_WEI in .env")
	}

	client := dialClient(cfg.RPCURL)
	defer client.Close()

	auth, _ := signer(cfg)

	// 已部署的合约地址（请替换为你部署时输出的地址）
	contractAddr := common.HexToAddress("<PUT-YOUR-CONTRACT-ADDRESS>")

	erc20, err := token.NewMyToken(contractAddr, client)
	if err != nil {
		log.Fatalf("bind err: %v", err)
	}

	to := common.HexToAddress(cfg.ToAddress)
	tx, err := erc20.Transfer(auth, to, cfg.TransferWei)
	if err != nil {
		log.Fatalf("transfer failed: %v", err)
	}
	fmt.Printf("💸 Transfer submitted: %s\n", tx.Hash().Hex())

	// 等待上链（可选）
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("wait mined: %v", err)
	}
	fmt.Println("✅ Transfer mined.")
}
