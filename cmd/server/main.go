package main

import (
	log "github.com/sirupsen/logrus"
)

func Run() error {
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Errorln(err)
	}
}
