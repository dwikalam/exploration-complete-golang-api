package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/config"
	"github.com/dwikalam/ecommerce-service/internal/app/db"
	"github.com/dwikalam/ecommerce-service/internal/app/handlers"
	"github.com/dwikalam/ecommerce-service/internal/app/loggers"
	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
	"github.com/dwikalam/ecommerce-service/internal/app/routes"
	"github.com/dwikalam/ecommerce-service/internal/app/services"
	"github.com/dwikalam/ecommerce-service/internal/app/transaction"
)

func Run(
	ctx context.Context,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var (
		defaultLogger loggers.Default
		cfg           config.Config
		psqlDB        db.Psql

		txManager transaction.SQLTransactionManager

		testRepository repositories.Test
		userRepository repositories.User

		testService services.Test
		authService services.Auth

		testHandler handlers.Test
		authHandler handlers.Auth

		err error
	)

	defaultLogger = loggers.NewDefault(stdout, stderr)

	cfg, err = config.New()
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Databases
	psqlDB, err = db.NewPsql(cfg.Db.PsqlURL)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}
	_, err = psqlDB.CheckHealth(ctx)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Transaction Manager
	txManager, err = transaction.NewManager(&psqlDB)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Repositories
	testRepository, err = repositories.NewTest(&psqlDB)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}
	userRepository, err = repositories.NewUser(&psqlDB)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Services
	testService, err = services.NewTest(&txManager, &testRepository)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}
	authService, err = services.NewAuth(&txManager, &userRepository)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	// Handlers
	testHandler, err = handlers.NewTest(&defaultLogger, &testService)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}
	authHandler, err = handlers.NewAuth(&defaultLogger, &authService)
	if err != nil {
		defaultLogger.Error(err.Error())
		return err
	}

	srvMux := routes.NewHttpHandler(
		&defaultLogger,
		&testHandler,
		&authHandler,
	)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      http.TimeoutHandler(srvMux, cfg.Server.HandlerTimeout, cfg.Server.TimeoutMessage),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	var (
		wg sync.WaitGroup

		listenAndServe = func() {
			defaultLogger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))

			err := httpServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				defaultLogger.Error(fmt.Sprintf("error listening and serving %s: %s", httpServer.Addr, err))

				return
			}
		}

		shutdownGracefully = func(ctx context.Context, wg *sync.WaitGroup) {
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
	)

	go listenAndServe()

	wg.Add(1)
	go shutdownGracefully(ctx, &wg)

	wg.Wait()

	return nil
}
