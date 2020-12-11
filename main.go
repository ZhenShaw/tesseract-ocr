package main

import (
	"flag"
	"os"

	"github.com/zhenshaw/tesseract-ocr/server"
)

const defaultPort = "8080"
const defaultToken = ""

func main() {

	flag.Parse()

	port := defaultPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	token := defaultToken
	if os.Getenv("TOKEN") != "" {
		token = os.Getenv("TOKEN")
	}

	srv := &server.Srv{
		Port:  port,
		Token: token,
	}
	srv.Run()
}
