package main

import (
	"github.com/c12s/meridian/model"
	"github.com/c12s/meridian/service"
	"github.com/c12s/meridian/storage/etcd"
	"log"
	"time"
)

func main() {
	conf, err := model.ConfigFile()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := etcd.New(conf, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	service.Run(db, conf)
}
