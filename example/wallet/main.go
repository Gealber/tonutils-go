package main

import (
	"context"
	"log"
	"strings"

	"github.com/Gealber/tonutils-go/address"
	"github.com/Gealber/tonutils-go/liteclient"
	"github.com/Gealber/tonutils-go/tlb"
	"github.com/Gealber/tonutils-go/ton"
	"github.com/Gealber/tonutils-go/ton/wallet"
)

func main() {
	client := liteclient.NewConnectionPool()

	// connect to mainnet lite server
	err := client.AddConnection(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	api := ton.NewAPIClient(client)
	// bound all requests to single ton node
	ctx := client.StickyContext(context.Background())

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	words := strings.Split("birth pattern then forest walnut then phrase walnut fan pumpkin pattern then cluster blossom verify then forest velvet pond fiction pattern collect then then", " ")

	w, err := wallet.FromSeed(api, words, wallet.V3)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.Address())

	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("CurrentMasterchainInfo err:", err.Error())
		return
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	if balance.NanoTON().Uint64() >= 3000000 {
		addr := address.MustParseAddr("EQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqB2N")

		log.Println("sending transaction and waiting for confirmation...")

		// if destination wallet is not initialized you should use TransferNoBounce
		// regular Transfer has bounce flag, and TONs may be returned.

		// err = w.TransferNoBounce(ctx, addr, tlb.MustFromTON("0.003"),
		err = w.Transfer(ctx, addr, tlb.MustFromTON("0.003"),
			"Hello from tonutils-go!", true)
		if err != nil {
			log.Fatalln("Transfer err:", err.Error())
			return
		}

		// update chain info
		block, err = api.CurrentMasterchainInfo(ctx)
		if err != nil {
			log.Fatalln("CurrentMasterchainInfo err:", err.Error())
			return
		}

		balance, err = w.GetBalance(ctx, block)
		if err != nil {
			log.Fatalln("GetBalance err:", err.Error())
			return
		}

		log.Println("transaction sent, balance left:", balance.TON())

		return
	}

	log.Println("not enough balance:", balance.TON())
}
