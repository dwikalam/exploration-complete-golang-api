package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dwikalam/ecommerce-service/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Run(
		ctx,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
