package main

import (
	"flag"
	"fmt"
	"log"

	cache "github.com/AniketDubey199/Cache-Proxy/Internal/Cache"
	proxy "github.com/AniketDubey199/Cache-Proxy/Internal/Proxy"
	"github.com/gofiber/fiber/v3"
)

func main() {
	port := flag.String("port", "", "Port to run proxy")
	origin := flag.String("origin", "", "origin server URL")

	clearCache := flag.Bool("clear-cache", false, "Clear cache")

	flag.Parse()

	cacheStore := cache.NewCache()

	if *clearCache {
		cacheStore.Clear()
		fmt.Println("Cache Cleared Succefully")
		return
	}

	if *port == "" || *origin == "" {
		log.Fatal("Usage: caching-proxy --port <number> --origin <url>")
	}

	app := fiber.New()

	handler := &proxy.ProxyHandler{
		Origin: *origin,
		Cache:  cacheStore,
	}

	app.All("/*", handler.Caching)

	log.Fatal(app.Listen(":" + *port))

}
