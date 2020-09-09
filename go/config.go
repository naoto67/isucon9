package main

import (
	"database/sql"
	"log"
)

var (
	ConfigCache = make(map[string]string)
)

func FlushConfig() {
	ConfigCache = make(map[string]string)
}

func getConfigByName(name string) (string, error) {
	if v, ok := ConfigCache[name]; ok {
		return v, nil
	}
	config := Config{}

	err := dbx.Get(&config, "SELECT * FROM `configs` WHERE `name` = ?", name)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		log.Print(err)
		return "", err
	}
	ConfigCache[config.Name] = config.Val
	return config.Val, err
}
