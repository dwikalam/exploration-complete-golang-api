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

	baselogger, err := loggers.NewBaseLogger(stdout, stderr)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	// Databases
	psqlDB, err := db.NewPostgresqlDB(&baselogger, cfg.PsqlURL)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}
	_, err = psqlDB.Health(ctx)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	// Repositories
	testRepository, err := repositories.NewTestRepo(&baselogger, &psqlDB)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	// Services
	testService, err := services.NewTestService(&baselogger, &testRepository)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	// Handlers
	testHandler, err := handlers.NewTestHandler(&baselogger, &testService)
	if err != nil {
		baselogger.Error(err.Error())
		return err
	}

	httpServer := &http.Server{
		Addr: net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
		Handler: routes.NewHttpHandler(
			&baselogger,
			&testHandler,
		),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 2,
	}

	listenAndServe := func() {
		baselogger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			baselogger.Error(fmt.Sprintf("error listening and serving %s: %s", httpServer.Addr, err))

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
			baselogger.Error(fmt.Sprintf("shutting down http server: %s", err))
		}

		baselogger.Warn("server shutdown")

		err = psqlDB.Disconnect()
		if err != nil {
			baselogger.Error(fmt.Sprintf("closing database: %s", err))
		}

		baselogger.Warn("database closed")
	}

	wg.Add(1)
	go shutdownGracefully()

	wg.Wait()

	return nil
}
