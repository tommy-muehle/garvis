package main

import (
	"flag"
	"net/http"

	"github.com/tommy-muehle/garvis/internal"
	"github.com/tommy-muehle/garvis/internal/configuration"
)

var (
	configPath = flag.String("config", "", "Path to configuration file.")
)

func main() {
	config := configuration.DefaultConfig()
	if *configPath != "" {
		config = configuration.MustLoadFromFile(*configPath)
	} else {
		config = configuration.FromEnvironment()
	}

	resolver := internal.NewResolver(config)
	webhook := resolver.ResolveWebhook()

	server := resolver.ResolveServer()
	server.Health("/health")
	server.AddHandler("/favicon.ico", http.NotFoundHandler())
	server.AddHandler("/payload", webhook.Handler())

	server.ListenAndServe()
}
