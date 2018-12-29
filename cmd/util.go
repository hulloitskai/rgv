package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: no .env file was found in the local "+
			"directory.")
	}
}

func buildLogger() (*zap.SugaredLogger, error) {
	var (
		raw *zap.Logger
		err error
	)
	if os.Getenv("GO_ENV") == "development" {
		raw, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.Encoding = "console"
		raw, err = cfg.Build()
	}
	if err != nil {
		return nil, err
	}
	return raw.Sugar(), nil
}
