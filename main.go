package main

import (
	"cheetah/framework"
	"log"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr: ":8888",
	}
	log.Fatal(server.ListenAndServe())
}
