package main

import (
	"fmt"
	"os"

	"github.com/Dzikuri/openidea-segokuning/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)

	// Call function connection database
	db, err := config.NewDatabase()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres connection error: %s", err.Error()))
		return
	}

	echo := config.NewEcho(&logger)

	config.Bootstrap(&config.BootstrapConfig{
		DB:     db,
		App:    echo,
		Logger: &logger,
	})

	echo.Logger.Fatal(echo.Start(":8080"))
}
