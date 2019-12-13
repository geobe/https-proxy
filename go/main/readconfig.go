package main

import (
	"fmt"
	"github.com/geobe/https-proxy/go/controller"
	"github.com/geobe/https-proxy/go/model"
)

func main() {
	controller.Setup("")
	model.Users = controller.ReadUsers()
	fmt.Printf("Users: %v\n", model.Users)
}
