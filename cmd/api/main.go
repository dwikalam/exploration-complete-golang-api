package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/config"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.HostServer, config.PortServer),
		Handler: nil,
	}

	startListenAndServe := func() {
		fmt.Fprintf(stdout, "listening on %s\n", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving %s\n", httpServer.Addr)
		}
	}
	go startListenAndServe()

	var wg sync.WaitGroup

	shutdownGracefully := func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error while shutting down http server: %s\n", err)
		}
	}

	wg.Add(1)
	go shutdownGracefully()

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()

	if err := run(
		ctx,
		os.Args,
		os.Getenv,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
