package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v10"

	"github.com/mihailozarinschi/ports-app/internal/adapters/memory"
	"github.com/mihailozarinschi/ports-app/internal/port"
	"github.com/mihailozarinschi/ports-app/internal/server"
)

type config struct {
	ServerAddr string `env:"SERVER_ADDR" envDefault:":3000"`

	PortRepositoryStorage string `env:"PORT_REPOSITORY_STORAGE" envDefault:"memory"`
}

func main() {
	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()
	mainWG := sync.WaitGroup{}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error parsing configs: %s", err)
	}

	var portsRepo port.Repository
	switch cfg.PortRepositoryStorage {
	case "memory":
		portsRepo = memory.NewPortRepository()
	default:
		log.Fatalf("unsupported storage %q for the port repository", cfg.PortRepositoryStorage)
	}

	// Start HTTP Server
	mainWG.Add(1)
	go func() {
		defer mainWG.Done()
		defer cancelMain()

		// Init server handlers
		srvHandler := server.NewHandler(portsRepo)

		// Init serve mux
		srvMux := server.NewServeMux(srvHandler)

		// TODO: Add logging, metrics, tracing, authorization, etc. as HTTP middlewares

		// Init server
		srv := http.Server{
			Addr:        cfg.ServerAddr,
			Handler:     srvMux,
			BaseContext: func(_ net.Listener) context.Context { return mainCtx },
		}

		// Handle graceful shutdown
		mainWG.Add(1)
		go func() {
			defer mainWG.Done()
			<-mainCtx.Done()
			err := srv.Shutdown(mainCtx)
			if err != nil {
				log.Printf("error shutting down http server: %s", err)
			}
		}()

		// Start listening for requests
		log.Printf("HTTP server listening on %q ...", cfg.ServerAddr)
		err := srv.ListenAndServe() // this blocks here, on purpose
		if err != nil {
			log.Println(err.Error())
		}
	}()

	// Handle os.Signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigCh
		log.Printf("Received os.Signal %q, shutting down...", s.String())
		cancelMain()
	}()

	// Wait for all processes to finish before exiting
	mainWG.Wait()
	log.Println("portsd process exiting now...")
}
