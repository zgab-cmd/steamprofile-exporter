# steamprofile-exporter

A small Prometheus exporter for Steam profile status.

## Build

On Windows, from the project root:

```powershell
go build -o steamprofile-exporter.exe .
```

This produces `steamprofile-exporter.exe`.

## Configuration

The exporter reads configuration from environment variables:

- `STEAM_API_KEY` — your Steam Web API key
- `STEAM_ID` — your Steam ID or vanity profile name (e.g. `YOUR_STEAM_ID`)
- `LISTEN_ADDR` — optional listen address, default `:9805`

The exporter already supports Steam vanity names via `ResolveVanityURL`, so users can provide either a numeric SteamID or a custom URL name.

## Metrics

The exporter exposes the following Prometheus metrics on `/metrics`:

- `steam_up` - gauge, `1` when the Steam API call succeeded, `0` when it failed.
- `steam_scrape_duration_seconds` - gauge, how long the Steam API call took in seconds.
- `steam_in_game` - gauge, `1` when the player appears to be in a game, otherwise `0`.
- `steam_current_game_info{game="..."}` - info-style gauge with the current game name as a label. This metric is only present when Steam returns a game name.
- `steam_persona_state` - gauge, raw numeric persona state from Steam (`0` to `7`).
- `steam_persona_state_info{state="..."}` - info-style gauge with the human-readable state name, such as `online`, `offline`, `away`, or `busy`.
- `steam_last_logoff_timestamp` - gauge, the Unix timestamp of the last logoff sent by Steam.


### Example `.env`

Copy `.env.example` to `.env` and fill in values:

```dotenv
STEAM_API_KEY=YOUR_API_KEY_HERE
STEAM_ID=YOUR_STEAM_ID
LISTEN_ADDR=:9805
```

## Run locally

From PowerShell:

```powershell
# load env vars into the current shell
Get-Content .env | ForEach-Object {
  if ($_ -match '^\s*([^#=]+)=(.*)$') {
    $name = $matches[1].Trim()
    $value = $matches[2].Trim().Trim('"')
    Set-Item -Path Env:$name -Value $value
  }
}

# run the exporter
.
\steamprofile-exporter.exe
```

Then open:

```powershell
curl.exe http://localhost:9805/metrics
```

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

For production use, install it as a Windows service. First set machine-level environment variables:

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

