package main

import (
	"net/url"
	"os"

	"github.com/frozzare/go-cfg"
)

type MyConfig struct {
	Name string
	URL  url.URL
}

func main() {
	mcfg := MyConfig{
		Name: "Fredrik",
	}

	// mcfg.Name => Fredrik

	// Create and fill a configuration object.
	c, _ := cfg.New(&mcfg, cfg.WithData(map[string]interface{}{
		"Name": "Test",
		"URL":  "https://github.com",
	}))

	// mcfg.Name => Test
	// mcfg.URL => https://github.com
	os.Setenv("TEST_NAME", "Go")

	// Extend existing configuration object with more values.
	c.Extend(cfg.WithEnvironment(map[string]string{
		"Name": "TEST_NAME",
	}))

	// mcfg.Name => Go
}
