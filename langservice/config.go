package langservice

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

const (
	LanguageServer Protocol = "LS"
	TSServer       Protocol = "TS"
)

var lsConfigFile string

var RootPathFinders = map[string]func() string{
	"go": findRootPathGo,
}

type Protocol string

type Config struct {
	Commands []string `json:"commands"`
	Protocol Protocol `json:"protocol"`
}

func LoadConfigFile(filename string) (map[string]Config, error) {
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	result := map[string]Config{}
	if err := json.NewDecoder(fin).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func LoadDefaultConfigFile() (map[string]Config, error) {
	if !pflag.Parsed() {
		return nil, errors.New("flags are not parsed yet")
	}
	return LoadConfigFile(lsConfigFile)
}

func findRootPathGo() string {
	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		return ""
	}
	path := filepath.ToSlash(filepath.Join(gopath, "src"))
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func DefaultRootPathFinder() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func init() {
	pflag.StringVar(&lsConfigFile, "ls-config", "", "Specify the filename containing JSON data of language servers.")
}