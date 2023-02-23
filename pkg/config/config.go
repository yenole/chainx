package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultEnv    = "CX_ENV"
	defaultPrefix = "CX_"
)

var (
	Data   = flag.String("data", "", "data folder")
	Listen = flag.String("listen", ":8080", "api server listen")
)

func init() {
	if v := os.Getenv(defaultEnv); v != "" {
		byts, _ := ioutil.ReadFile(v)
		list := strings.Split(strings.TrimSpace(string(byts)), "\n")
		for _, line := range list {
			if strings.HasPrefix(line, defaultPrefix) && strings.Contains(line, "=") {
				kv := strings.Split(line, "=")
				flag.Set(strings.ToLower(strings.TrimPrefix(kv[0], defaultPrefix)), strings.TrimSpace(kv[1]))
			}
		}
	}

	for _, line := range os.Environ() {
		if strings.HasPrefix(line, defaultPrefix) && strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			flag.Set(strings.ToLower(strings.TrimPrefix(kv[0], defaultPrefix)), strings.TrimSpace(kv[1]))
		}
	}
}

func JoinPath(args ...string) string {
	return filepath.Join(*Data, filepath.Join(args...))
}
