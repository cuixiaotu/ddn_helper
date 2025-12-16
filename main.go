package main

import (
	"log"
	"os"
	"tb_helper/ui"
)

func main() {
	logger := log.New(os.Stdout, "ui_", log.LstdFlags)
	home := ui.NewHomePage(logger)
	home.Box()
}
