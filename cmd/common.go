package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type Cfg struct {
	RPCURL           string
	PrivateKeyHex    string
	ChainID          *big.Int
	TokenName        string
	TokenSymbol      string
	InitialSupplyWei *big.Int
	ToAddress        string
	TransferWei      *big.Int
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func loadCfg() *Cfg {
	_ = godotenv.Load() // 允许 .env

	chainID := new(big.Int)
	chainID, ok := chainID.SetString(mustEnv("CHAIN_ID"), 10)
	if !ok {
		log.Fatal("invalid CHAIN_ID")
	}

	initSupply := new(big.Int)
	initSupply, ok = initSupply.SetString(mustEnv("INITIAL_SUPPLY_WEI"), 10)
	if !ok {
		log.Fatal("invalid INITIAL_SUPPLY_WEI")
	}

	transferWei := new(big.Int)
	if v := os.Getenv("TRANSFER_AMOUNT_WEI"); v != "" {
		if _, ok := transferWei.SetString(v, 10); !ok {
			log.Fatal("invalid TRANSFER_AMOUNT_WEI")
		}
	}

	return &Cfg{
		RPCURL:           mustEnv("RPC_URL"),
		PrivateKeyHex:    mustEnv("PRIVATE_KEY"),
		ChainID:          chainID,
		TokenName:        mustEnv("TOKEN_NAME"),
		TokenSymbol:      mustEnv("TOKEN_SYMBOL"),
		InitialSupplyWei: initSupply,
		ToAddress:        os.Getenv("TO_ADDRESS"),
		TransferWei:      transferWei,
	}
}

func dialClient(rpc string) *ethclient.Client {
	c, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatalf("dial RPC failed: %v", err)
	}
	return c
}

func signer(cfg *Cfg) (*bind.TransactOpts, *ecdsa.PrivateKey) {
	keyHex := strings.TrimPrefix(cfg.PrivateKeyHex, "0x")
	priv, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		log.Fatalf("bad PRIVATE_KEY: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(priv, cfg.ChainID)
	if err != nil {
		log.Fatalf("transactor err: %v", err)
	}
	// 让 geth 自动估算费用；如需手动设置可在这里指定 GasTipCap/GasFeeCap
	auth.Context = context.Background()
	return auth, priv
}
