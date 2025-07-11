package main

import (
	"github.com/jiujuan/go-star/internal/app"
	"github.com/jiujuan/go-star/internal/router"
)

func main() {
	app.Bootstrap(
		app.Modules,
		router.Module,
	)
}