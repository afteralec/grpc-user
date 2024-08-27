package main

import (
	"context"
	"fmt"
	"os"

	"github.com/afteralec/grpc-user/server"
	"github.com/spf13/viper"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: Move shared validators into a Petrichormud package

// TODO: Structured logging
func main() {
	ctx := context.Background()
	if err := server.Run(ctx, newConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "err from server: %s", err)
		os.Exit(1)
	}
}

func newConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigType("toml")
	config.AddConfigPath("/run/secrets")

	config.SetConfigName("root_username")
	config.ReadInConfig()
	config.SetConfigName("root_passphrase")
	config.MergeInConfig()

	return config
}
