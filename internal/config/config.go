package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
  DBURL string `json:"db_url"`
  CurrentUserName string `json:"current_user_name"`
}

func getFilePath () (string, error) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return "", err
  }
  file := filepath.Join(homeDir, configFileName)
  if err != nil {
    return "", err
  }
  return file, nil
}

func (cfg *Config) SetUser (userName string) error {
  cfg.CurrentUserName = userName
  err := write(*cfg)
  if err != nil {
    return  err
  }
  return nil
}

func write (cfg Config) error {
  file, err := getFilePath()
  if err != nil {
    return err
  }
  data, err := os.Create(file)
  if err != nil {
    return err
  }
  defer data.Close()
  encoder := json.NewEncoder(data)
  if err := encoder.Encode(&cfg); err != nil {
    return err
  }
  return nil
}

func Read () (Config, error) {
  file, err := getFilePath()
  if err != nil {
    return Config{}, fmt.Errorf("Error while filePath to json file: %w", err)
  }
  data, err := os.Open(file)
  if err != nil {
    return Config{}, fmt.Errorf("Error while opening the file to read: %w", err)
  }
  defer data.Close()
  var cfg Config
  decoder := json.NewDecoder(data)
  if err := decoder.Decode(&cfg); err != nil {
    return Config{}, err
  }
  return cfg, nil
}


