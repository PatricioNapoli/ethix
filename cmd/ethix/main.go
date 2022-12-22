package main

import (
	"context"
	"github.com/PatricioNapoli/ethix/pkg/indexer"
	"github.com/PatricioNapoli/ethix/pkg/server"
	"github.com/PatricioNapoli/ethix/pkg/store"
	"github.com/PatricioNapoli/ethix/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Printf("launching Ethix")

	addrDb := store.NewGobDatabase[string]("addresses.gob")
	txnsDb := store.NewGobDatabase[[]indexer.Transaction]("transactions.gob")

	go utils.DoEvery(10*time.Second, func() {
		log.Printf("flushing databases to disk")

		addrDb.Persist()
		txnsDb.Persist()
	})

	ctx, cancelCtx := context.WithCancel(context.Background())

	ctx = context.WithValue(ctx, "indexer", indexer.New("https://cloudflare-eth.com/", &addrDb, &txnsDb))

	sv := server.NewHTTPServer(":8080", ctx, cancelCtx)

	server.RegisterRoutes(sv)

	go sv.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	<-stop

	log.Println("shutting down Ethix")

	sv.Server.Close()
}
