package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// deve inicializar um context com cancelamento para receber sinais de término
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// carrega as configurações
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed o open config file, %s", err)
	}
	server := NewServer(config)
	errChan := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		errChan <- server.Run(ctx)
	})
	select {
	case <-ctx.Done():
		log.Println("received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Printf("server stopped, %s", err)
		}
	}
	wg.Wait()
}
