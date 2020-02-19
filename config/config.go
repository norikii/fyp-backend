package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type DBConfig struct {
	DatabaseName string `json:"database_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Url string `json:"url"`
}

func (dc *DBConfig) GetConfig(fileName string) (*DBConfig, error) {
	conf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &DBConfig{}, fmt.Errorf("could not read file %s: %v", fileName, err)
	}

	err = json.Unmarshal(conf, dc)
	if err != nil {
		return &DBConfig{}, fmt.Errorf("could not unmarshal file %s: %v", fileName, err)
	}

	return dc, nil
}


