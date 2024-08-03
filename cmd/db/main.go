package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	upArg   = "-up"
	downArg = "-down"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Println(err)

		os.Exit(1)
	}
}

func run(args []string) error {
	const (
		migrationsDir string = "file://cmd/db/migrations"
	)

	var (
		arg string
		cfg config.EnvConfig
		m   *migrate.Migrate

		err error
	)

	cfg, err = config.NewEnvConfig()
	if err != nil {
		return fmt.Errorf("creating env config failed: %v", err)
	}

	m, err = migrate.New(
		migrationsDir,
		cfg.GetDbPsqlDSN(),
	)
	if err != nil {
		return fmt.Errorf("creating migrate failed: %v", err)
	}

	arg, err = validatedArgs(args)
	if err != nil {
		return fmt.Errorf("arg validation failed: %v", err)
	}

	switch arg {
	case upArg:
		return m.Up()
	case downArg:
		return m.Down()
	default:
		return errors.New("no migrate commands executed")
	}
}

func validatedArgs(args []string) (string, error) {
	var (
		validArgs  string
		argsLength = len(args)
	)

	if argsLength != 2 {
		return "", errors.New("invalid args length. only 1 given arg is allowed")
	}

	validArgs = args[1]

	if validArgs != upArg && validArgs != downArg {
		return "", fmt.Errorf("invalid arg. allowed arg is either %s or %s", upArg, downArg)
	}

	return validArgs, nil
}
