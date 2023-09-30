package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/env"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		CleanEnv         bool   `flag:"clean" default:"false" description:"Do not pass current environment to child process"`
		EncryptionMethod string `flag:"encryption" default:"openssl-md5" description:"Encryption method used for encrypted env-file (Available: gpg-symmetric, openssl-md5, openssl-sha256)"`
		EnvFile          string `flag:"env-file" default:".env" description:"Location of the environment file"`
		LogLevel         string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		PasswordFile     string `flag:"password-file" default:"" description:"Read encryption key from file"`
		Password         string `flag:"password,p" default:"" env:"PASSWORD" description:"Password to decrypt environment file"`
		VersionAndExit   bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return fmt.Errorf("parsing cli options: %w", err)
	}

	l, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log-level: %w", err)
	}
	log.SetLevel(l)

	return nil
}

func main() {
	var err error

	if err = initApp(); err != nil {
		log.WithError(err).Fatal("intitializing app")
	}

	if cfg.VersionAndExit {
		fmt.Printf("envrun %s\n", version) //nolint:forbidigo
		os.Exit(0)
	}

	if cfg.Password == "" && cfg.PasswordFile != "" {
		if _, err := os.Stat(cfg.PasswordFile); err == nil {
			data, err := os.ReadFile(cfg.PasswordFile)
			if err != nil {
				log.WithError(err).Fatal("Unable to read password from file")
			}
			cfg.Password = strings.TrimSpace(string(data))
		}
	}

	dec, err := decryptMethodFromName(cfg.EncryptionMethod)
	if err != nil {
		log.WithError(err).Fatal("Could not load decrypt method")
	}

	pairs, err := loadEnvFromFile(cfg.EnvFile, cfg.Password, dec)
	if err != nil {
		log.WithError(err).Fatal("Could not load env file")
	}

	childenv := env.ListToMap(os.Environ())
	if cfg.CleanEnv {
		childenv = map[string]string{}
	}

	for k, v := range pairs {
		childenv[k] = v
	}

	if len(rconfig.Args()) < 2 { //nolint:gomnd
		log.Fatal("No command specified")
	}

	c := exec.Command(rconfig.Args()[1], rconfig.Args()[2:]...) //#nosec:G204 // Intended to run cmd from input
	c.Env = env.MapToList(childenv)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err = c.Run()

	switch err.(type) {
	case nil:
		log.Info("Process exitted with code 0")
		os.Exit(0)
	case *exec.ExitError:
		log.Error("Unclean exit with exit-code != 0")
		os.Exit(1)
	default:
		log.WithError(err).Error("An unknown error occurred")
		os.Exit(2) //nolint:gomnd
	}
}

func loadEnvFromFile(filename, passphrase string, decrypt decryptMethod) (map[string]string, error) {
	body, err := os.ReadFile(filename) //#nosec:G304 // Intended to read a variable env file
	if err != nil {
		return nil, fmt.Errorf("reading env-file: %w", err)
	}

	if passphrase != "" {
		if body, err = decrypt(body, passphrase); err != nil {
			return nil, fmt.Errorf("decrypting env-file: %w", err)
		}
	}

	return env.ListToMap(strings.Split(string(body), "\n")), nil
}
