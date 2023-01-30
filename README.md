<div align="center" width="100%">
    <img src="frontend/static/favicon.png" width="128" />
</div>

<div align="center" width="100%">
    <h2>UpSnap</h2>
    <p>A simple wake on lan app written with SvelteKit, Go, PocketBase and nmap.</p>
    <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/stars/seriousm4x/upsnap" /></a>
    <a target="_blank" href="https://github.com/seriousm4x/UpSnap/releases"><img src="https://img.shields.io/github/downloads/seriousm4x/upsnap/total" /></a>
    <a target="_blank" href="https://github.com/seriousm4x/UpSnap/releases"><img src="https://img.shields.io/github/v/release/seriousm4x/upsnap?display_name=tag" /></a>
    <a target="_blank" href="https://github.com/seriousm4x/UpSnap/actions"><img src="https://github.com/seriousm4x/upsnap/actions/workflows/deploy.yml/badge.svg?event=push" /></a>
    <a target="_blank" href="https://github.com/seriousm4x/UpSnap/commits/master"><img src="https://img.shields.io/github/last-commit/seriousm4x/upsnap" /></a>
    <h3><a href="https://github.com/seriousm4x/UpSnap/releases/tag/3.0.0">v3 is here! 🤩 Check it out</a></h3>
</div>

## ✨ Features

- Dashboard to wake up devices with one click
- Set timed wake and shutdown events via cron
- Add custom ports to devices which will be pinged
- Discover devices by scanning network
- Dark/light or system prefered color scheme
- [Docker images](https://github.com/seriousm4x/UpSnap/pkgs/container/upsnap) for amd64, arm64, arm/v7

## 📸 Screenshots

| Dark                           | Light                           |
| ------------------------------ | ------------------------------- |
| ![](/assets/home_dark.png)     | ![](/assets/home_light.png)     |
| ![](/assets/device_dark.png)   | ![](/assets/device_light.png)   |
| ![](/assets/settings_dark.png) | ![](/assets/settings_light.png) |

## 🚀 Run the binary

Just download the latest binary from the [release page](https://github.com/seriousm4x/UpSnap/releases) and run it.

Run `./upsnap serve --http=0.0.0.0:8090`. For more options check `./upsnap --help`.

If you need network discovery, make sure to have nmap installed and run upsnap as root/admin.

## 🐳 Run in docker

Alternatively use the [docker-compose](docker-compose.yml) example. See the comments in the file for customization.

### Reverse Proxy

**Caddy example**

```
upsnap.example.com {
    reverse_proxy localhost:8090
}
```

## Help developing

Fork this branch and clone it.

1. Start backend

```sh
cd backend
go mod tidy
go run main.go serve
```

Log in to [http://localhost:8090/\_/](http://localhost:8090/_/), create an admin user and add some devices.

2. Start frontend

```sh
cd frontend
pnpm i
pnpm run dev
```

Open up [http://localhost:5173/](http://localhost:5173/)
