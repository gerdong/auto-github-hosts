package main

import (
	"github.com/gerdong/auto-github-hosts/config"
	"github.com/gerdong/auto-github-hosts/domain"
)

func main2() {
	config.Init()
	domain.Cron()
}
