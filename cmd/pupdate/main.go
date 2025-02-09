package main

import (
	"github.com/langet/pupdate/internal/routes"
)

func main() {
	newRouter := routes.NewRouter()
	newRouter.Run(":8080")
}
