package main

import (
	"github.com/c12s/meridian/model"
	"github.com/c12s/meridian/service"
	"github.com/c12s/meridian/storage/etcd"
	"github.com/c12s/meridian/storage/redis"
	"log"
	"time"
)

func main() {
	conf, err := model.ConfigFile()
	if err != nil {
		log.Fatal(err)
		return
	}

	cache, err := redis.New(conf.Cache)
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := etcd.New(conf, cache, 10*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}

	service.Run(db, conf)
}
