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
	"github.com/dwikalam/ecommerce-service/internal/app/loggers"
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

	var (
		defaultLogger loggers.DefaultLogger
		cfg           config.Config
		psqlDB        db.PostgresqlDB

		testRepository repositories.TestRepository
		testService    services.TestService
		testHandler    handlers.TestHandler

		err error
	)

	defaultLogger = loggers.NewDefaultLogger(stdout, stderr)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	cfg, err = config.New()
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Databases
	psqlDB, err = db.NewPostgresqlDB(&defaultLogger, cfg.PsqlURL)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}
	_, err = psqlDB.CheckHealth(ctx)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Repositories
	testRepository, err = repositories.NewTestRepo(&defaultLogger, &psqlDB)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Services
	testService, err = services.NewTestService(&defaultLogger, &testRepository)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Handlers
	testHandler, err = handlers.NewTestHandler(&defaultLogger, &testService)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
		Handler:      routes.NewHttpHandler(&defaultLogger, &testHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 2,
	}

	listenAndServe := func() {
		defaultLogger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			defaultLogger.Error(fmt.Sprintf("error listening and serving %s: %s", httpServer.Addr, err))

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
			defaultLogger.Error(fmt.Sprintf("shutting down http server: %s", err))
		}

		defaultLogger.Warn("server shutdown")

		err = psqlDB.Disconnect()
		if err != nil {
			defaultLogger.Error(fmt.Sprintf("closing database: %s", err))
		}

		defaultLogger.Warn("database closed")
	}

	wg.Add(1)
	go shutdownGracefully()

	wg.Wait()

	return nil
}
