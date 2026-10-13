package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// SteamCollector implements prometheus.Collector and calls the Steam API on each scrape.
type SteamCollector struct {
	apiKey  string
	steamID string

	// descriptors for metrics
	steamUpDesc                  *prometheus.Desc
	steamScrapeDurationDesc      *prometheus.Desc
	steamInGameDesc              *prometheus.Desc
	steamCurrentGameInfoDesc     *prometheus.Desc
	steamPersonaStateDesc        *prometheus.Desc
	steamPersonaStateInfoDesc    *prometheus.Desc
	steamLastLogoffTimestampDesc *prometheus.Desc
	steamPlaytimeTotalHoursDesc  *prometheus.Desc
}

// NewSteamCollector creates a new collector instance.
func NewSteamCollector(apiKey, steamID string) *SteamCollector {
	return &SteamCollector{
		apiKey:                       apiKey,
		steamID:                      steamID,
		steamUpDesc:                  prometheus.NewDesc("steam_up", "1 if all Steam API calls succeeded", nil, nil),
		steamScrapeDurationDesc:      prometheus.NewDesc("steam_scrape_duration_seconds", "Total duration of Steam API requests in seconds", nil, nil),
		steamInGameDesc:              prometheus.NewDesc("steam_in_game", "1 if the player is currently in a game", nil, nil),
		steamCurrentGameInfoDesc:     prometheus.NewDesc("steam_current_game_info", "Info metric for the current game (label=game)", []string{"game"}, nil),
		steamPersonaStateDesc:        prometheus.NewDesc("steam_persona_state", "Numeric persona state from Steam (0-7)", nil, nil),
		steamPersonaStateInfoDesc:    prometheus.NewDesc("steam_persona_state_info", "Info metric for persona state (label=state)", []string{"state"}, nil),
		steamLastLogoffTimestampDesc: prometheus.NewDesc("steam_last_logoff_timestamp", "Unix timestamp of last logoff", nil, nil),
		steamPlaytimeTotalHoursDesc:  prometheus.NewDesc("steam_playtime_total_hours", "Total playtime across owned games in hours", nil, nil),
	}
}

// Describe sends the descriptors for each metric.
func (s *SteamCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.steamUpDesc
	ch <- s.steamScrapeDurationDesc
	ch <- s.steamInGameDesc
	ch <- s.steamCurrentGameInfoDesc
	ch <- s.steamPersonaStateDesc
	ch <- s.steamPersonaStateInfoDesc
	ch <- s.steamLastLogoffTimestampDesc
	ch <- s.steamPlaytimeTotalHoursDesc
}

var personaStateNames = map[int]string{
	0: "offline",
	1: "online",
	2: "busy",
	3: "away",
	4: "snooze",
	5: "looking_to_trade",
	6: "looking_to_play",
	7: "max",
}

// Collect is called on each scrape — it performs Steam API requests synchronously,
// then emits metrics. On failure, only steam_up=0 and scrape duration are emitted.
func (s *SteamCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	ps, err := FetchPlayerSummary(s.apiKey, s.steamID)
	var totalPlaytimeHours float64
	if err == nil {
		totalPlaytimeHours, err = FetchPlayerPlaytimeTotalHours(s.apiKey, s.steamID)
	}
	duration := time.Since(start).Seconds()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(s.steamUpDesc, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(s.steamScrapeDurationDesc, prometheus.GaugeValue, duration)
		return
	}

	// success
	ch <- prometheus.MustNewConstMetric(s.steamUpDesc, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(s.steamScrapeDurationDesc, prometheus.GaugeValue, duration)

	if ps.GameExtraInfo != "" {
		ch <- prometheus.MustNewConstMetric(s.steamInGameDesc, prometheus.GaugeValue, 1)
		ch <- prometheus.MustNewConstMetric(s.steamCurrentGameInfoDesc, prometheus.GaugeValue, 1, ps.GameExtraInfo)
	} else {
		ch <- prometheus.MustNewConstMetric(s.steamInGameDesc, prometheus.GaugeValue, 0)
	}

	ch <- prometheus.MustNewConstMetric(s.steamPersonaStateDesc, prometheus.GaugeValue, float64(ps.PersonaState))
	stateName, ok := personaStateNames[ps.PersonaState]
	if !ok {
		stateName = "unknown"
	}
	ch <- prometheus.MustNewConstMetric(s.steamPersonaStateInfoDesc, prometheus.GaugeValue, 1, stateName)

	ch <- prometheus.MustNewConstMetric(s.steamLastLogoffTimestampDesc, prometheus.GaugeValue, float64(ps.LastLogoff))
	ch <- prometheus.MustNewConstMetric(s.steamPlaytimeTotalHoursDesc, prometheus.GaugeValue, totalPlaytimeHours)
}
