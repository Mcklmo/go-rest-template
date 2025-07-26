package config

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port                    string
	Address                 string
	PrivateKey              any
	PublicKey               any
	DatabaseURL             string
	RootPath                string
	MuscleCSVPath           string
	AssetPath               string
	PublicKeyLocalFilePath  string
	PrivateKeyLocalFilePath string
	DropDatabase            bool
}

func NewConfig() *Config {
	c := &Config{
		Port:                    "8080",
		Address:                 "0.0.0.0",
		DatabaseURL:             "./workout_tracker.db",
		MuscleCSVPath:           "muscle.csv",
		AssetPath:               "src/asset",
		PublicKeyLocalFilePath:  "public.pem",
		PrivateKeyLocalFilePath: "private.pem",
		DropDatabase:            true,
	}

	flagSet := flag.NewFlagSet("server", flag.ExitOnError)
	flagSet.StringVar(&c.RootPath, "rootPath", "..", "Define the path to the root of the server. This should be '.' if run from the root and '..' if run using the debugger.")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	c.PublicKey = parseKey(loadKey(c.RootPath+"/"+c.PublicKeyLocalFilePath), x509.ParsePKIXPublicKey)
	c.PrivateKey = parseKey(loadKey(c.RootPath+"/"+c.PrivateKeyLocalFilePath), x509.ParsePKCS8PrivateKey)

	return c
}

func parseKey(rawKey []byte, parseFunc func(der []byte) (key any, err error)) (parsedKey any) {
	var (
		block *pem.Block
		err   error
	)

	if block, _ = pem.Decode(rawKey); block == nil {
		panic("failed to parse PEM block")
	}

	if parsedKey, err = parseFunc(block.Bytes); err != nil {
		panic(fmt.Sprintf("failed to parse key: %s", err))
	}

	return parsedKey
}

func loadKey(path string) []byte {
	rawKey, err := os.ReadFile(path) //nolint:gosec // path is not user input
	if err != nil {
		panic(err)
	}

	return rawKey
}
