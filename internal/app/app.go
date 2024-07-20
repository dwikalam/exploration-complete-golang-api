package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/config"
	"github.com/dwikalam/ecommerce-service/internal/app/db"
	"github.com/dwikalam/ecommerce-service/internal/app/handlers"
	"github.com/dwikalam/ecommerce-service/internal/app/routes"
	"github.com/dwikalam/ecommerce-service/internal/app/types/customerr"
)

func Run(
	ctx context.Context,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := db.Initialize(); err != nil && !errors.Is(err, &customerr.DatabaseAlreadyConnectedError{}) {
		return err
	}

	httpServer := &http.Server{
		Addr: net.JoinHostPort(config.ServerHost, config.ServerPort),
		Handler: routes.NewHttpHandler(
			handlers.NewTestHandler(),
		),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	listenAndServe := func() {
		fmt.Fprintf(stdout, "listening on %s\n", httpServer.Addr)

		err := httpServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving %s\n", httpServer.Addr)

			return
		}
	}
	go listenAndServe()

	var wg sync.WaitGroup

	shutdownGracefully := func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}

	wg.Add(1)
	go shutdownGracefully()

	wg.Wait()

	return nil
}
