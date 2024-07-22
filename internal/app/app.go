package app

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

	"github.com/dwikalam/ecommerce-service/internal/app/config"
	"github.com/dwikalam/ecommerce-service/internal/app/db"
	"github.com/dwikalam/ecommerce-service/internal/app/handlers"
	"github.com/dwikalam/ecommerce-service/internal/app/logger"
	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
	"github.com/dwikalam/ecommerce-service/internal/app/routes"
	"github.com/dwikalam/ecommerce-service/internal/app/services"
)

func Run(
	ctx context.Context,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger, err := logger.NewBaseLogger(stdout, stderr)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Databases
	psqlDB, err := db.NewPostgresqlDB(&logger, cfg.PsqlURL)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	_, err = psqlDB.Health(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Repositories
	testRepository, err := repositories.NewTestRepo(&logger, &psqlDB)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Services
	testService, err := services.NewTestService(&logger, &testRepository)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Handlers
	testHandler, err := handlers.NewTestHandler(&logger, &testService)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	httpServer := &http.Server{
		Addr: net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
		Handler: routes.NewHttpHandler(
			&logger,
			&testHandler,
		),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 2,
	}

	listenAndServe := func() {
		logger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("error listening and serving %s: %s", httpServer.Addr, err))

			return
		}
	}
	go listenAndServe()

	var wg sync.WaitGroup

	shutdownGracefully := func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error(fmt.Sprintf("shutting down http server: %s", err))
		}

		logger.Warn("server shutdown")

		err = psqlDB.Disconnect()
		if err != nil {
			logger.Error(fmt.Sprintf("closing database: %s", err))
		}

		logger.Warn("database closed")
	}

	wg.Add(1)
	go shutdownGracefully()

	wg.Wait()

	return nil
}
