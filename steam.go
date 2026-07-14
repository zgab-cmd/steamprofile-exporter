package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// PlayerSummary contains a subset of fields returned by Steam's GetPlayerSummaries.
type PlayerSummary struct {
	PersonaState  int    `json:"personastate"`
	GameExtraInfo string `json:"gameextrainfo"`
	LastLogoff    int64  `json:"lastlogoff"`
}

// responseEnvelope models the outer response from the Steam API.
type responseEnvelope struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
}

// FetchPlayerSummary synchronously calls the Steam Web API and returns the first player summary.
// This is a simple HTTP GET and JSON decode. The function is deliberately synchronous because
// the exporter will call it once per metrics scrape.
func FetchPlayerSummary(apiKey, steamID string) (*PlayerSummary, error) {
	// Steam Web API requires a 64-bit SteamID for GetPlayerSummaries. If the
	// provided steamID looks like a vanity name (non-numeric), try to resolve
	// it to a 64-bit id using ResolveVanityURL.
	if _, err := strconv.ParseUint(steamID, 10, 64); err != nil {
		resolved, err := ResolveVanityURL(apiKey, steamID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve vanity Steam ID: %w", err)
		}
		steamID = resolved
	}

	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", apiKey, steamID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned status %d", resp.StatusCode)
	}

	var env responseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}

	if len(env.Response.Players) == 0 {
		return nil, fmt.Errorf("no player data in response")
	}
	return &env.Response.Players[0], nil
}

// ResolveVanityURL converts a vanity name (custom profile name) to a 64-bit SteamID.
// It calls the Steam Web API ResolveVanityURL method and returns the numeric steamid.
func ResolveVanityURL(apiKey, vanity string) (string, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", apiKey, vanity)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("resolve vanity API returned status %d", resp.StatusCode)
	}

	var env struct {
		Response struct {
			SteamID string `json:"steamid"`
			Success int    `json:"success"`
			Message string `json:"message"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", err
	}
	if env.Response.Success != 1 {
		return "", fmt.Errorf("vanity resolution failed: %s", env.Response.Message)
	}
	if env.Response.SteamID == "" {
		return "", fmt.Errorf("no steamid returned from vanity resolution")
	}
	return env.Response.SteamID, nil
}
