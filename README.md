<div align="center" width="100%">
    <img src="app/frontend/public/favicon.png" width="128" />
</div>

<div align="center" width="100%">
    <h2>UpSnap</h2>
    <p>A simple wake on lan app written with Svelte, Django, Django-Channels (websockets), Celery, Redis and nmap.</p>
    <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/stars/seriousm4x/upsnap" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/pulls/seriousm4x/upsnap" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/v/seriousm4x/upsnap/latest?label=docker%20image%20ver." /></a> <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/last-commit/seriousm4x/upsnap" /></a>
</div>

## ✨ Features

* Dashboard to wake up devices with 1 click
* Set timed wake and shutdown events via cron
* Add custom ports to devices which will be scanned
* Discover devices by scanning network
* Notifications on status changes
* Devices only get pinged when there are 1 or more visitors
* Dark/light or system prefered color scheme
* [Docker images](https://hub.docker.com/r/seriousm4x/upsnap) for amd64, arm64, arm/v7

## 📸 Screenshots

| Dark                 | Light                 |
| -------------------- | --------------------- |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/index-dark.png.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/index-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/device-settings-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/device-settings-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/settings-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/settings-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/add-device-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/add-device-light.png) |

## 🐳 Run your own instance

There are 3 example docker-compose files to choose from. The simplest is [docker-compose-sqlite.yml](docker-compose-sqlite.yml).

The website will be available at [localhost:8000](http://localhost:8000). If you run it on a different pc, it will be `http://<your-ip>:8000`. You can change the port in the docker-compose file.

### Databases

Upsnap supports 3 different databases. Postgres, MySQL and SQLite. If you already have an existing database you want to use, delete the database container from the compose file. Always make sure to set the correct database type environment variable, e.g. DB_TYPE=mysql

## 📝 Other infos

* The app container needs to run in host network mode to send the wakeonlan command on your local network. Therefore all other containers also need to run in host network mode. I don't like it but there is no way around.
