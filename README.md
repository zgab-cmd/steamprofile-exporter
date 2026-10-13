# steamprofile-exporter

A small Prometheus exporter for Steam playtime and presence info.

95% vibecoded using Raptor mini (~5 credits in total), but works perfectly well. Contributions are welcome!

## Quick start

Grab a binary from the [releases page](https://github.com/zgab-cmd/steamprofile-exporter/releases) (Windows `.exe`, Linux amd64/arm64), or use the Docker image:

```bash
docker run --rm -p 9805:9805 \
  -e STEAM_API_KEY='YOUR_API_KEY' \
  -e STEAM_ID='YOUR_STEAM_ID_OR_VANITY' \
  ghcr.io/zgab-cmd/steamprofile-exporter:latest
```

Then check <http://localhost:9805/metrics>.

## Configuration

Set these environment variables (see [.env.example](.env.example)):

| Variable | Description |
|---|---|
| `STEAM_API_KEY` | Your Steam Web API key, get one at <https://steamcommunity.com/dev/apikey> |
| `STEAM_ID` | Numeric SteamID64 or vanity name (the part after `steamcommunity.com/id/`) |
| `LISTEN_ADDR` | Optional listen address, default `:9805` |

The Steam API is only queried when `/metrics` is scraped; nothing is polled in the background.

## Metrics

Exposed on `/metrics`:

- `steam_up` - `1` if all Steam API calls succeeded, `0` if any failed.
- `steam_scrape_duration_seconds` - total duration of the Steam API requests.
- `steam_in_game` - `1` when the player appears to be in a game, otherwise `0`.
- `steam_current_game_info{game="..."}` - current game name as a label; only present when Steam returns one.
- `steam_persona_state` - raw persona state from Steam (`0` to `7`).
- `steam_persona_state_info{state="..."}` - human-readable state (`online`, `offline`, `away`, `busy`, ...).
- `steam_last_logoff_timestamp` - Unix timestamp of the last logoff.
- `steam_playtime_total_hours` - total playtime across owned games, including played free games, in hours.

Note: if your Steam profile is private, the persona state will always be `0`. Steam may return no owned-game data for private game details, in which case total playtime is reported as `0`.

## Build from source

Requires Go.

```bash
go build -o steamprofile-exporter .
```

Cross-compile for Linux with `GOOS=linux GOARCH=amd64` (or `arm64`). Build the image with `docker build -t steamprofile-exporter .`.

## Run

Linux / macOS:

```bash
export STEAM_API_KEY='YOUR_API_KEY' STEAM_ID='YOUR_STEAM_ID'
./steamprofile-exporter
```

Windows (PowerShell):

```powershell
$env:STEAM_API_KEY='YOUR_API_KEY'
$env:STEAM_ID='YOUR_STEAM_ID'
.\steamprofile-exporter.exe
```

### Kubernetes / k3s

Create the secret, then apply the manifest in [k8s/](k8s/steamprofile-exporter-deployment.yaml):

```bash
kubectl create secret generic steamprofile-exporter-secrets \
  --from-literal=steam_api_key='YOUR_API_KEY' \
  --from-literal=steam_id='YOUR_STEAM_ID_OR_VANITY'

kubectl apply -f k8s/steamprofile-exporter-deployment.yaml
```

Expose it with an ingress or LoadBalancer if you need access from outside the cluster.

### Windows service (optional)

Set machine-level environment variables, then register the service (run as Administrator):

```powershell
[Environment]::SetEnvironmentVariable('STEAM_API_KEY','YOUR_API_KEY','Machine')
[Environment]::SetEnvironmentVariable('STEAM_ID','YOUR_STEAM_ID','Machine')

New-Service -Name SteamProfileExporter -BinaryPathName 'C:\path\to\steamprofile-exporter.exe' -DisplayName 'Steam Profile Exporter' -StartupType Automatic
Start-Service SteamProfileExporter
```

## License

[MIT](LICENSE)
