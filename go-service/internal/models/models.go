package models

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ListenPort          int
	BackendBaseUrl      string
	BackendUserName     string
	BackendUserPassword string
}

func GetConfig() (Config, error) {

	fnCheckEnv := func(envName, envDesc string) (string, error) {
		value := os.Getenv(envName)
		if strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("no config %v specified", envDesc)
		}

		return value, nil
	}

	cfg := Config{}

	value, err := fnCheckEnv("PORT", "port")
	if err != nil {
		return Config{}, err
	}

	cfg.ListenPort, err = strconv.Atoi(value)
	if err != nil {
		return Config{}, err
	}

	value, err = fnCheckEnv("BACKEND_BASE_URL", "base url")
	if err != nil {
		return Config{}, err
	}
	cfg.BackendBaseUrl = value

	value, err = fnCheckEnv("BACKEND_USER_NAME", "backend user name")
	if err != nil {
		return Config{}, err
	}
	cfg.BackendUserName = value

	value, err = fnCheckEnv("BACKEND_USER_PASSWORD", "backend user password")
	if err != nil {
		return Config{}, err
	}
	cfg.BackendUserPassword = value

	return cfg, nil
}
