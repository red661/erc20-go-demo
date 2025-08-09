package main

import (
	"fmt"
	"log"

	"erc20-go-demo/internal/token"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	cfg := loadCfg()
	client := dialClient(cfg.RPCURL)
	defer client.Close()

	auth, priv := signer(cfg)

	// 部署合约
	addr, tx, instance, err := token.DeployMyToken(
		auth, client,
		cfg.TokenName,
		cfg.TokenSymbol,
		cfg.InitialSupplyWei,
	)
	if err != nil {
		log.Fatalf("deploy failed: %v", err)
	}
	fmt.Printf("🚀 Deploying MyToken... tx=%s\n", tx.Hash().Hex())
	fmt.Printf("📜 Contract address: %s\n", addr.Hex())

	// 查询：总供应 & 部署者余额
	total, _ := instance.TotalSupply(nil)
	pub := priv.PublicKey
	from := crypto.PubkeyToAddress(pub)
	bal, _ := instance.BalanceOf(nil, from)
	name, _ := instance.Name(nil)
	sym, _ := instance.Symbol(nil)
	dec, _ := instance.Decimals(nil)

	fmt.Printf("Token: %s (%s) decimals=%d\n", name, sym, dec)
	fmt.Printf("TotalSupply: %s\n", total.String())
	fmt.Printf("Deployer balance: %s\n", bal.String())

	fmt.Println("✅ Deployed.")
}
