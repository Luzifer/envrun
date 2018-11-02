package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig"
)

var (
	cfg = struct {
		CleanEnv         bool   `flag:"clean" default:"false" description:"Do not pass current environment to child process"`
		EncryptionMethod string `flag:"encryption" default:"openssl-md5" description:"Encryption method used for encrypted env-file (Available: gpg-symmetric, openssl-md5, openssl-sha256)"`
		EnvFile          string `flag:"env-file" default:".env" description:"Location of the environment file"`
		LogLevel         string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		PasswordFile     string `flag:"password-file" default:"" description:"Read encryption key from file"`
		Password         string `flag:"password,p" default:"" env:"PASSWORD" description:"Password to decrypt environment file"`
		Silent           bool   `flag:"q" default:"false" description:"Suppress informational messages from envrun (DEPRECATED, use --log-level=warn)"`
		VersionAndExit   bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("envrun %s\n", version)
		os.Exit(0)
	}

	if cfg.Silent && cfg.LogLevel == "info" {
		// Migration of deprecated flag
		cfg.LogLevel = "warn"
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func envListToMap(list []string) map[string]string {
	out := map[string]string{}
	for _, entry := range list {
		if len(entry) == 0 || entry[0] == '#' {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		out[parts[0]] = parts[1]
	}
	return out
}

func envMapToList(envMap map[string]string) []string {
	out := []string{}
	for k, v := range envMap {
		out = append(out, k+"="+v)
	}
	return out
}

func main() {
	if cfg.Password == "" && cfg.PasswordFile != "" {
		if _, err := os.Stat(cfg.PasswordFile); err == nil {
			data, err := ioutil.ReadFile(cfg.PasswordFile)
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

	var childenv = envListToMap(os.Environ())
	if cfg.CleanEnv {
		childenv = map[string]string{}
	}

	for k, v := range pairs {
		childenv[k] = v
	}

	c := exec.Command(rconfig.Args()[1], rconfig.Args()[2:]...)
	c.Env = envMapToList(childenv)
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
		os.Exit(2)
	}
}

func loadEnvFromFile(filename, passphrase string, decrypt decryptMethod) (map[string]string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not read env-file: %s", err)
	}

	if passphrase != "" {
		if body, err = decrypt(body, passphrase); err != nil {
			return nil, fmt.Errorf("Could not decrypt env-file: %s", err)
		}
	}

	return envListToMap(strings.Split(string(body), "\n")), nil
}
