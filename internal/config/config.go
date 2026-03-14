package config

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("error occurred while fetching user home directory: %v", err)
		return "", err
	}
	return filepath.Join(homeDir, "xsh"), nil
}

func CheckConfigDir() bool {
	configDir, err := GetConfigDir()
	if err != nil {
		return false
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Debugf("Config directory does not exist at: %s", configDir)
		return false
	}

	log.Debugf("Config directory exists at: %s", configDir)
	return true
}

func InitConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Errorf("error occurred while creating config directory: %v", err)
		return err
	}

	log.Debugf("Config directory created at: %s", configDir)
	return nil
}
