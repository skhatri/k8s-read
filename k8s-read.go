package main

import (
	"flag"
	"fmt"
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/k8s-read/controller"
	"net/http"
	"os"
	"strconv"
)

func parseArguments(args []string) string {
	app := flag.NewFlagSet("serve", flag.ExitOnError)
	var port = 6100
	httpPortFromEnv := os.Getenv("HTTP_PORT")
	if httpPortFromEnv != "" {
		port, _ = strconv.Atoi(httpPortFromEnv)
	}
	var address = "0.0.0.0"
	app.StringVar(&address, "host", address, "Host Interface to listen on")
	app.IntVar(&port, "port", port, "Web port to bind to")
	app.Parse(args[1:])
	if app.Parsed() {
		if port < 1024 || port > 65535 {
			panic("invalid port passed. please provide one between 1024 and 65535")
		}
	}

	return fmt.Sprintf("%s:%d", address, port)
}

func main() {
	var args []string
	if len(os.Args) < 2 {
		args = []string{
			"serve",
		}
	} else {
		args = os.Args[1:]
	}
	var command = args[0]

	switch command {
	case "serve":
		address := parseArguments(args)
		r := router.NewHttpRouterBuilder().
			WithOptions(router.HttpRouterOptions{
			LogRequest:  true,
			LogFunction: func(values ...interface{}) {
				fmt.Println(values)
			},
		}).Configure(func (configurer router.ApiConfigurer){
			controller.Configure(configurer)
		}).Build()
		fmt.Printf("Listening on %s\n", address)
		http.ListenAndServe(address, r)
	default:
		fmt.Printf("command %s is not supported\n", command)
	}

}
