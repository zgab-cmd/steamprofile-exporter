# steamprofile-exporter

A small Prometheus exporter for exporting playtime infos from Steam.
95% vibecoded using Raptor mini, only using ~5 credits in total

## Build

On Windows, from the project root:

```powershell
go build -o steamprofile-exporter.exe .
```

This produces `steamprofile-exporter.exe`.

## Configuration

The exporter reads configuration from environment variables:

- `STEAM_API_KEY` — your Steam Web API key
- `STEAM_ID` — your Steam ID or vanity profile name (the numbers or string after steamcommunity.com/id/...)
- `LISTEN_ADDR` — optional listen address, default `:9805`

Note: get your API key from here: https://steamcommunity.com/dev
The exporter already supports Steam vanity names via `ResolveVanityURL`, so users can provide either a numeric SteamID or a custom URL name.
The Steam API is only queried upon scraping the /metrics page, so it does not poll Steam continuously in the background.

## Metrics

The exporter exposes the following Prometheus metrics on `/metrics`:

- `steam_up` - gauge, `1` when the Steam API call succeeded, `0` when it failed.
- `steam_scrape_duration_seconds` - gauge, how long the Steam API call took in seconds.
- `steam_in_game` - gauge, `1` when the player appears to be in a game, otherwise `0`.
- `steam_current_game_info{game="..."}` - info-style gauge with the current game name as a label. This metric is only present when Steam returns a game name.
- `steam_persona_state` - gauge, raw numeric persona state from Steam (`0` to `7`).
- `steam_persona_state_info{state="..."}` - info-style gauge with the human-readable state name, such as `online`, `offline`, `away`, or `busy`.
- `steam_last_logoff_timestamp` - gauge, the Unix timestamp of the last logoff sent by Steam.

Note: if your profile is set to private, your  state will always show up as 0

### Example `.env`

Copy `.env.example` to `.env` and fill in values:

```dotenv
STEAM_API_KEY=YOUR_API_KEY_HERE
STEAM_ID=YOUR_STEAM_ID_HERE
LISTEN_ADDR=:9805
```

## Run locally

Build the executable first:

```powershell
go build -o steamprofile-exporter.exe 
```

From PowerShell:

```powershell
# create .env from the example if needed

Copy-Item .env.example .env -ErrorAction SilentlyContinue


# load env vars into the current shell
Get-Content .env | ForEach-Object {
  if ($_ -match '^\s*([^#=]+)=(.*)$') {
    $name = $matches[1].Trim()
    $value = $matches[2].Trim().Trim('"')
    Set-Item -Path Env:$name -Value $value
  }
}

# run the exporter
.\steamprofile-exporter.exe
```

Then open:

```powershell
curl.exe http://localhost:9805/metrics
```

## Linux / Docker support

You can also run this exporter on Linux or in a container.

Build a Linux executable with:

```bash
go build -o steamprofile-exporter-linux .
```

Build a Docker image with:

```bash
docker build -t steamprofile-exporter:latest .
```

Run the container with:

```bash
docker run --rm -p 9805:9805 \
  -e STEAM_API_KEY='YOUR_API_KEY' \
  -e STEAM_ID='YOUR_STEAM_ID_OR_VANITY' \
  -e LISTEN_ADDR=':9805' \
  steamprofile-exporter:latest
```

### Kubernetes / k3s

A Kubernetes deployment manifest is available at `k8s/steamprofile-exporter-deployment.yaml`.

Create the secret for your credentials:

```bash
kubectl create secret generic steamprofile-exporter-secrets \
  --from-literal=steam_api_key='YOUR_API_KEY' \
  --from-literal=steam_id='YOUR_STEAM_ID_OR_VANITY'
```

Apply the deployment:

```bash
kubectl apply -f k8s/steamprofile-exporter-deployment.yaml
```

If you want the exporter available outside the cluster, expose it with an ingress or LoadBalancer.

Or just run the pre built .exe, but you still have to create an .env file under your C:\Users\youruser\ and load them in, using the above script

## Run on another Windows PC

Copy these files to the other machine:

- `steamprofile-exporter.exe`
- `steam.go` and `main.go` are not needed once the exe is built
- `.env.example` or a custom `.env` with the user's own key and Steam ID

Then start it with the user's own values:

```powershell
$env:STEAM_API_KEY='USER_API_KEY'
$env:STEAM_ID='USER_STEAM_ID_OR_VANITY'
$env:LISTEN_ADDR=':9805'
.
\steamprofile-exporter.exe
```

## Service installation (optional)

For consistent use, install it as a Windows service. First set machine-level environment variables:

```powershell
[Environment]::SetEnvironmentVariable('STEAM_API_KEY','USER_API_KEY','Machine')
[Environment]::SetEnvironmentVariable('STEAM_ID','USER_STEAM_ID_OR_VANITY','Machine')
[Environment]::SetEnvironmentVariable('LISTEN_ADDR',':9805','Machine')
```

Then create and start the service:

```powershell
New-Service -Name SteamProfileExporter -BinaryPathName 'C:\path\to\steamprofile-exporter.exe' -DisplayName 'Steam Profile Exporter' -StartupType Automatic
Start-Service SteamProfileExporter
```

