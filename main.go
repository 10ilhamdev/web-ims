package main

import (
	"ims/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
