package types

import (
	"os"
)

var (
	VAULT_ADDR = os.Getenv("VAULT_ADDR")
	USERNAME = os.Getenv("VAULT_USER")
	PASSWORD = os.Getenv("VAULT_PASSWORD")
	ENCRYPT_KEY = os.Getenv("ENCRYPT_KEY")
	PROJECT_NAME = os.Getenv("VAULT_PROJECT_NAME")
)

type Secrets struct {
	Secret map[string]string `yaml:"secret"`
}

