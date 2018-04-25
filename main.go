package main

import (
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/web"
	"github.com/InitialShape/blockchain/utils"
	"log"
	"net/http"
	"os"
	"flag"
)

func main() {
	var store blockchain.Store
	var peer blockchain.Peer

	keys := flag.Bool("generate_keys", false,
					  "Generates keys for the wallet and miner")
	flag.Parse()
	if *keys {
			utils.GenerateWallet()
	} else {
		// normal operation mode
		store = blockchain.Store{}
		store.Open(os.Args[1], &peer)
		peer = blockchain.Peer{"localhost", os.Args[2], store}
		go peer.Start()

		_, err := store.StoreGenesisBlock(20)
		if err != nil {
			log.Fatal(err)
		}

		r := web.Handlers(store)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Args[3]), r))
	}

}
