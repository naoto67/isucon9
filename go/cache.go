package main

import (
	"database/sql"
	"log"
)

var (
	ConfigMap = make(map[string]string)
)

func getConfigByName(name string) (string, error) {
	if val, ok := ConfigMap[name]; ok {
		return val, nil
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

	ConfigMap[name] = config.Val

	return config.Val, err
}
