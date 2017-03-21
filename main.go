package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	openssl "github.com/Luzifer/go-openssl"
	"github.com/Luzifer/rconfig"
)

var (
	version = "dev"
	cfg     = struct {
		EnvFile  string `flag:"env-file" default:".env" description:"Location of the environment file"`
		Silent   bool   `flag:"q" default:"false" description:"Suppress informational messages from envrun"`
		CleanEnv bool   `flag:"clean" default:"false" description:"Do not pass current environment to child process"`
		Password string `flag:"password,p" default:"" env:"PASSWORD" description:"Password to decrypt environment file"`
	}{}
)

func init() {
	rconfig.Parse(&cfg)
}

func infoLog(message string, args ...interface{}) {
	if !cfg.Silent {
		log.Printf(message, args...)
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
		log.Fatalf("Could not read env-file: %s", err)
	}

	if cfg.Password != "" {
		if body, err = openssl.New().DecryptString(cfg.Password, string(body)); err != nil {
			log.Fatalf("Could not decrypt env-file: %s", err)
		}
	}

	var childenv map[string]string
	if cfg.CleanEnv {
		childenv = map[string]string{}
	} else {
		childenv = envListToMap(os.Environ())
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
		infoLog("Process exitted with code 0")
		os.Exit(0)
	case *exec.ExitError:
		infoLog("Unclean exit with exit-code != 0")
		os.Exit(1)
	default:
		log.Printf("An unknown error ocurred: %s", err)
		os.Exit(2)
	}
}
