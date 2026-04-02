<div align="center" width="100%">
    <img src="frontend/static/gopher.svg" width="150" />
</div>

<div align="center" width="100%">
    <h2>UpSnap</h2>
    <p>A simple wake on lan web app written with SvelteKit, Go and PocketBase.</p>
    <div>
        <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/stars/seriousm4x/UpSnap?style=flat&label=Stars" /></a>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/pkgs/container/upsnap"><img src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fraw.githubusercontent.com%2Fipitio%2Fbackage%2Findex%2Fseriousm4x%2FUpSnap%2Fupsnap.json&query=downloads&label=ghcr.io%20pulls" /></a>
        <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/pulls/seriousm4x/upsnap?label=docker%20hub%20pulls" /></a>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/releases"><img src="https://img.shields.io/github/downloads/seriousm4x/upsnap/total?label=binary%20downloads" /></a>
    </div>
    <div>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/releases"><img src="https://img.shields.io/github/go-mod/go-version/seriousm4x/UpSnap?filename=backend/go.mod" /></a>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/releases"><img src="https://img.shields.io/github/v/release/seriousm4x/upsnap?display_name=tag&label=Latest%20release" /></a>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/actions"><img src="https://github.com/seriousm4x/upsnap/actions/workflows/release.yml/badge.svg?event=push" /></a>
        <a target="_blank" href="https://github.com/seriousm4x/UpSnap/commits/master"><img src="https://img.shields.io/github/last-commit/seriousm4x/upsnap" /></a>
    </div>
</div>

> [!NOTE]
>
> #### UpSnap is, and always will be, free and open source software.
>
> If someone is asking you to pay money for access to UpSnap binaries, source code, or licenses, you are being scammed.
>
> The official and only trusted source for UpSnap is this repository (and its linked releases).
> Do not pay third parties for something that is provided here for free.

## ✨ Features

- 🚀 One-Click Device Wake-Up Dashboard
- ⏰ Timed Events via Cron for Automation
- 🔌 Ping Any Port You Choose
- 🔍 Discover Devices with Network Scanning (nmap required)
- ❎️ Shutdown Devices with a Custom Command
- 👤 Secured User Management
- 🌐 i18n support for [these](/frontend/translations) languages
- 🎨 35 Themes
- 🐳 [Docker images](https://github.com/seriousm4x/UpSnap/pkgs/container/upsnap) for amd64, arm64, arm/v7, arm/v6
- 🏠 Self-Hostable

## 📸 Screenshots

| Silk                               | Dim                               |
| ---------------------------------- | --------------------------------- |
| ![](/assets/home-light.webp)       | ![](/assets/home-dark.webp)       |
| ![](/assets/account-light.webp)    | ![](/assets/account-dark.webp)    |
| ![](/assets/new-manual-light.webp) | ![](/assets/new-manual-dark.webp) |
| ![](/assets/new-scan-light.webp)   | ![](/assets/new-scan-dark.webp)   |
| ![](/assets/settings-light.webp)   | ![](/assets/settings-dark.webp)   |
| ![](/assets/users-light.webp)      | ![](/assets/users-dark.webp)      |

## 🚀 Run the binary

Just download the latest binary from the [release page](https://github.com/seriousm4x/UpSnap/releases) and run it.

```sudo ./upsnap serve --http=0.0.0.0:8090```

To run as non-root (Linux only), refer to the Wiki: [Use non-root user](https://github.com/seriousm4x/UpSnap/wiki/Use-non%E2%80%90root-user)

For more options check `./upsnap --help` or visit [PocketBase documentation](https://pocketbase.io/docs).

## 🐳 Run in Docker

Just pull and run the image: `docker run --network=host seriousm4x/upsnap:latest` or you can use the [docker-compose](docker-compose.yml) example. See the comments in the file for customization.

To run as non-root, refer to the Wiki: [Use non-root user](https://github.com/seriousm4x/UpSnap/wiki/Use-non%E2%80%90root-user)

### Change port (Docker)

If you want to change the port from 8090 to something else, change the following (5000 in this case):

```yml
environment:
  - UPSNAP_HTTP_LISTEN=0.0.0.0:5000
```
## ❎️ Shutting Down Devices

To shutdown devices, refer to the Wiki: [How to use shutdowns](https://github.com/seriousm4x/UpSnap/wiki/How-to-use-shutdowns)
## Reverse Proxy

**Caddy example**

```
upsnap.example.com {
    reverse_proxy localhost:8090
}
```

## Run in sub path

You can run UpSnap on a different path than `/`, e.g. `/upsnap-sub-path/`. To do this in caddy, set the following:

```
http://localhost:8091 {
    handle /upsnap-sub-path/* {
        uri strip_prefix /upsnap-sub-path
        reverse_proxy localhost:8090
    }
}
```

Or nginx:

```
http {
    server {
        listen 8091;
        server_name localhost;
        location /upsnap-sub-path/ {
            proxy_pass http://localhost:8090/;
            proxy_redirect off;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

Paths must end with a trailing `/`.

## 🐧 Install from the [AUR](https://aur.archlinux.org/packages/upsnap-bin)

```bash
yay -Sy upsnap-bin
```

## 🔒 User permissions

UpSnap offers unique access for each user, per device. While admins have all permissions, they can assign specific rights to users such as displaying/hiding a device, accessing device editing, deleting and powering devices on/off. See the last screenshot in the [📸 Screenshots section](#-screenshots).

## 🌍 Exposing to the open web

Although UpSnap has user authorisation, it is **not recommended to expose it to the open web** and make it accessible by everyone!

**Reason**: The shutdown device command is basically a command piped to #sh (possibly as root user). If anyone gains unauthorized access and can abuse this API route in any way, the attacker has access to a (root) shell on your local network.

**Recommended**: If you need access from outside your network, please use a VPN. Wireguard or OpenVPN is your way to go.

## 🌐 Help translating

UpSnap is available in the following languages so far:

- 🇸🇦 **Arabic (Saudi Arabia)** (ar-SA)
- 🇧🇬 **Bulgarian (Bulgaria)** (bg-BG)
- 🇨🇿 **Czech (Czech republic)** (cs-CZ)
- 🇩🇪 **German (Germany)** (de-DE)
- 🇺🇸 **English (United States)** (en-US)
- 🇪🇸 **Spanish (Spain)** (es-ES)
- 🇫🇷 **French (France)** (fr-FR)
- 🇮🇳 **Hindi (India)** (hi-IN)
- 🇮🇩 **Bahasa (Indonesia)** (id-ID)
- 🇮🇹 **Italian (Italy)** (it-IT)
- 🇯🇵 **Japanese (Japan)** (ja-JP)
- 🇰🇷 **Korean (Republic of Korea)** (ko-KR)
- 🇳🇱 **Dutch (Netherlands)** (nl-NL)
- 🇳🇴 **Norwegian (Norway)** (nb-NO)
- 🇵🇱 **Polish (Poland)** (pl-PL)
- 🇵🇹 **Portuguese (Portugal)** (pt-PT)
- 🇷🇺 **Russian (Russia)** (ru-RU)
- 🇺🇦 **Ukrainian (Ukrain)** (uk-UA)
- 🇻🇳 **Vietnamese (Vietnam)** (vi-VN)
- 🇹🇼 **Chinese (Taiwan)** (zh-TW)
- 🇨🇳 **Chinese (China)** (zh-CN)

**If you want to contribute and help translating, check the wiki: [How to add languages](https://github.com/seriousm4x/UpSnap/wiki/How-to-add-languages)**

## 🔧 Help developing

Fork this branch and clone it.

1. Start backend

```sh
cd backend
go mod tidy
go run main.go serve
```

2. Start frontend

```sh
cd frontend
pnpm i
pnpm run dev
```

Open up [http://localhost:5173/](http://localhost:5173/), create an admin user and add some devices.

3. Optional: run the MAUI wrapper app

```sh
cd mobile
dotnet build -t:Run -f net10.0-windows10.0.19041.0
```

Then open `Settings` inside the app and point it to the frontend URL you want to use.

## 📱 Mobile app

This repository also includes a .NET MAUI app in [mobile](mobile) called **UpSnap Mobile**.

It is intentionally simple: the MAUI app does not reimplement the UI. It opens your existing UpSnap frontend inside a native WebView and gives you a small native settings screen where you can save the frontend URL on the device.

That makes it useful when:

- you already host UpSnap on your LAN or behind a reverse proxy
- you want a native app shell for Android, Windows, iOS or MacCatalyst
- you want to keep the web app as the single UI source of truth

### How it works

- The mobile app lives in [mobile/mobile.csproj](mobile/mobile.csproj)
- The saved frontend URL is configured inside the app and persisted locally on the device
- The main screen is a WebView that points to your existing UpSnap frontend
- Settings opens as a native page, so changing the URL does not depend on the web app being reachable

### Important note about URLs

The mobile app needs the **real URL that the device can reach**.

- `http://localhost:5173` only works if the frontend is running on the same machine as the app
- for a phone or tablet, use a LAN IP such as `http://192.168.1.50:5173` or your hosted URL such as `https://upsnap.example.com`
- for the Android emulator, `http://10.0.2.2:5173` maps back to the host machine

### Prerequisites

You will need:

- .NET 10 SDK
- .NET MAUI workloads for the target platforms you want to build
- Android SDK / emulator for Android builds
- Xcode and Apple signing setup for iOS or MacCatalyst builds

To inspect installed workloads:

```sh
dotnet workload list
```

If needed, install MAUI workloads:

```sh
dotnet workload install maui
```

### Run the mobile app in development

1. Make sure the UpSnap frontend is already running and reachable
2. Move into the mobile project

```sh
cd mobile
```

3. Build it once

```sh
dotnet build
```

4. Run a target platform

Windows:

```sh
dotnet build -t:Run -f net10.0-windows10.0.19041.0
```

Android:

```sh
dotnet build -t:Run -f net10.0-android
```

After the app opens:

1. Open `Settings`
2. Enter the full frontend URL
3. Save it
4. Return to the main screen and the WebView will load UpSnap

### Supported targets

The current MAUI project targets:

- Android: `net10.0-android`
- iOS: `net10.0-ios`
- MacCatalyst: `net10.0-maccatalyst`
- Windows: `net10.0-windows10.0.19041.0`

Application details:

- App name: `UpSnap Mobile`
- App ID: `com.upsnap.mobile`
- Windows packaging: unpackaged (`WindowsPackageType=None`)

### Publish release builds

All examples below assume you are already inside the `mobile` folder.

If you publish from **Visual Studio**, keep one practical rule in mind:

- use a separate publish or archive flow for each target platform
- do not reuse a Windows publish profile for Android, or the other way around

If Visual Studio ever reports a missing runtime pack such as `Microsoft.NETCore.App.Runtime.win-x64` while you are trying to publish Android, it usually means the publish profile or restore state is targeting the wrong platform. In that case, recreate the Android publish profile or use the Android archive flow directly.

#### Windows

This produces an unpackaged Windows release build:

```sh
dotnet restore -r win-x64
dotnet publish -c Release -f net10.0-windows10.0.19041.0 -r win-x64
```

If you want a self-contained Windows build:

```sh
dotnet publish -c Release -f net10.0-windows10.0.19041.0 -r win-x64 --self-contained true
```

Output is typically under:

```text
mobile/bin/Release/net10.0-windows10.0.19041.0/win-x64/publish/
```

#### Android

Release APK build:

```sh
dotnet publish -c Release -f net10.0-android -p:AndroidPackageFormat=apk
```

If you want an Android App Bundle (`.aab`) for store distribution:

```sh
dotnet publish -c Release -f net10.0-android -p:AndroidPackageFormat=aab
```

For Play Store or signed release packages, you will also need your own keystore and signing properties.

If you publish from **Visual Studio** for Android, the usual path is:

1. Select the Android target
2. Switch to `Release`
3. Use `Archive` / `Distribute` for the Android app
4. Choose APK or AAB in that Android-specific flow

If Visual Studio still complains about a Windows runtime pack while publishing Android, delete the old publish profile and create a fresh Android one. That error is usually not about Android itself, but about the wrong target or runtime being reused.

#### MacCatalyst

```sh
dotnet publish -c Release -f net10.0-maccatalyst
```

If you need both Intel and Apple Silicon runtimes, set runtime identifiers explicitly when packaging.

#### iOS

```sh
dotnet publish -c Release -f net10.0-ios
```

iOS release publishing requires Apple signing and, in normal practice, a Mac build environment.

### Mobile project structure

The most relevant files are:

- [mobile/mobile.csproj](mobile/mobile.csproj): MAUI targets, app identity and assets
- [mobile/MauiProgram.cs](mobile/MauiProgram.cs): dependency injection setup
- [mobile/Services/FrontendSettings.cs](mobile/Services/FrontendSettings.cs): persisted frontend URL storage
- [mobile/Pages/BrowserPage.xaml](mobile/Pages/BrowserPage.xaml): WebView host page
- [mobile/Pages/SettingsPage.xaml](mobile/Pages/SettingsPage.xaml): native settings page for the frontend URL

### Practical release advice

- Test with the exact URL your real device will use before publishing
- Prefer a stable LAN hostname or reverse-proxied HTTPS URL for daily use
- Keep the mobile app thin and let the web frontend remain the product UI
- If the web app layout changes, retest the long forms in the WebView before shipping a release

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=seriousm4x/UpSnap&type=Date&theme=dark)](https://star-history.com/#seriousm4x/UpSnap&Date)
