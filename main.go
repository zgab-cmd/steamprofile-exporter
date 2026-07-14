package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main reads environment variables, registers the collector, and starts HTTP server.
func main() {
	apiKey := os.Getenv("STEAM_API_KEY")
	steamID := os.Getenv("STEAM_ID")
	listen := os.Getenv("LISTEN_ADDR")
	if listen == "" {
		listen = ":9805"
	}

	if apiKey == "" || steamID == "" {
		log.Fatal("STEAM_API_KEY and STEAM_ID must be set in environment")
	}

	// If the user provided a vanity Steam ID (custom URL name), attempt to resolve
	// it here so the collector always holds a numeric 64-bit SteamID.
	if _, err := strconv.ParseUint(steamID, 10, 64); err != nil {
		resolved, err := ResolveVanityURL(apiKey, steamID)
		if err != nil {
			log.Fatalf("failed to resolve STEAM_ID %s: %v", steamID, err)
		}
		log.Printf("Resolved STEAM_ID %s -> %s", steamID, resolved)
		steamID = resolved
	}

	collector := NewSteamCollector(apiKey, steamID)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s, serving /metrics", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
