package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	Conf Config
)

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "default"
	}
	log.Printf("loading config by env:【%s】...\n", env)
	cp := fmt.Sprintf("./configs/%s.json", env)
	c, err := os.ReadFile(cp)
	if err != nil {
		log.Fatalf("load env config failed! error:%s\n", err.Error())
	}
	if err = json.Unmarshal(c, &Conf); err != nil {
		log.Fatalf("load env config failed! error:%s\n", err.Error())
	}
}
