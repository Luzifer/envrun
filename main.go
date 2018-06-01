package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	openssl "github.com/Luzifer/go-openssl"
	"github.com/Luzifer/rconfig"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		EnvFile        string `flag:"env-file" default:".env" description:"Location of the environment file"`
		Silent         bool   `flag:"q" default:"false" description:"Suppress informational messages from envrun (DEPRECATED, use --log-level=warn)"`
		CleanEnv       bool   `flag:"clean" default:"false" description:"Do not pass current environment to child process"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		Password       string `flag:"password,p" default:"" env:"PASSWORD" description:"Password to decrypt environment file"`
		PasswordFile   string `flag:"password-file" default:"" description:"Read encryption key from file"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
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
	body, err := ioutil.ReadFile(cfg.EnvFile)
	if err != nil {
		log.WithError(err).Fatal("Could not read env-file")
	}

	if cfg.Password == "" && cfg.PasswordFile != "" {
		if _, err := os.Stat(cfg.PasswordFile); err == nil {
			data, err := ioutil.ReadFile(cfg.PasswordFile)
			if err != nil {
				log.WithError(err).Fatal("Unable to read password from file")
			}
			cfg.Password = string(data)
		}
	}

	if cfg.Password != "" {
		if body, err = openssl.New().DecryptString(cfg.Password, string(body)); err != nil {
			log.WithError(err).Fatal("Could not decrypt env-file")
		}
	}

	var childenv = envListToMap(os.Environ())
	if cfg.CleanEnv {
		childenv = map[string]string{}
	}

	pairs := envListToMap(strings.Split(string(body), "\n"))
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
		log.WithError(err).Error("An unknown error ocurred")
		os.Exit(2)
	}
}
