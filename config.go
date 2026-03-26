package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Nome do arquivo de configuração
const configFileName = "config.json"

// Configuração da aplicação.
type Config struct {
	ServerAddress string                  `json:"server_address"`
	Environments  map[string]*Environment `json:"environments"`
}

// Configuração dos ambientes onde a aplicação irá se conectar.
type Environment struct {
	StsAddress   string `json:"sts_address"`
	ApiAddress   string `json:"api_address"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Cria uma instancia das configurações da aplicação.
func NewConfig() *Config {
	return &Config{
		ServerAddress: "localhost:8080",
		Environments: map[string]*Environment{
			"dev": {
				StsAddress:   "https://sts.com",
				ApiAddress:   "https://api.com",
				ClientId:     "xxx",
				ClientSecret: "yyy",
			},
		},
	}
}

// Cria um arquivo com a configuração padrão.
func NewConfigFile() (config *Config, err error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = "./"
	}
	config = NewConfig()
	f, err := os.OpenFile(filepath.Join(dir, configFileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0744)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	err = enc.Encode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Carrega o arquivo de configuração ou cria um.
func LoadConfig() (config *Config, err error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = "./"
	}
	f, err := os.OpenFile(filepath.Join(dir, configFileName), os.O_RDONLY, 0744)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewConfigFile()
		}
	}
	defer f.Close()
	config = &Config{}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Salva o arquivo de configuração.
func (p *Config) Save() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = "./"
	}
	f, err := os.OpenFile(filepath.Join(dir, configFileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0744)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	err = enc.Encode(p)
	if err != nil {
		return err
	}
	return nil
}
