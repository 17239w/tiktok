package main

import (
	"fmt"
	"tiktok/config"
	"tiktok/router"
)

func main() {
	r := router.Init()
	err := r.Run(fmt.Sprintf(":%d", config.Global.Port)) // listen and serve on "localhost:1116" (for windows "localhost:8080")
	if err != nil {
		return
	}
}
