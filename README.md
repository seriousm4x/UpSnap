<div align="center" width="100%">
    <img src="app/frontend/public/favicon.png" width="128" />
</div>

<div align="center" width="100%">
    <h2>UpSnap</h2>
    <p>A simple wake on lan app written with Django, Django-Channels (websockets), Celery, Redis and nmap.</p>
    <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/stars/seriousm4x/upsnap" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/pulls/seriousm4x/upsnap" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/v/seriousm4x/upsnap/latest?label=docker%20image%20ver." /></a> <a target="_blank" href="https://github.com/seriousm4x/upsnap"><img src="https://img.shields.io/github/last-commit/seriousm4x/upsnap" /></a>
</div>

## ✨ Features

* Dashboard to wake up devices with 1 click
* Set timed wake events via cron
* Scan devices for open ports you assigned
* Discover devices by scanning network
* Visualization of pings
* Notifications on status changes
* Devices only get pinged when there are 1 or more visitors
* Dark/light mode via preferes-color-scheme
* [Docker images](https://hub.docker.com/r/seriousm4x/upsnap) for amd64, arm64, arm/v7

## 📸 Screenshots

| Dark                 | Light                 |
| -------------------- | --------------------- |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/front-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/front-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/settings-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/settings-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/schedule-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/upsnap/master/assets/schedule-light.png) |

## 🐳 Run your own instance

There are 3 example docker-compose files to choose from. The simplest is [docker-compose-sqlite.yml](docker-compose-sqlite.yml).

The website will be available at [localhost:8000](http://localhost:8000). If you run it on a different pc, it will be `http://<your-ip>:8000`. You can change the port in the docker-compose file.

### Databases

Upsnap supports 3 different databases. Postgres, MySQL and SQLite. If you already have an existing database you want to use, delete the database container from the compose file. Always make sure to set the correct database type environment variable, e.g. DB_TYPE=mysql

## 🔧 Available Env Vars

| env var | type | info |
|---------|------|------|
| DJANGO_SUPERUSER_USER | Str | Username for /admin backend |
| DJANGO_SUPERUSER_PASSWORD | Str | Password for /admin backend |
| DJANGO_SECRET_KEY | Str | You can create your own by running `docker exec upsnap_django bash -c "python manage.py shell -c 'from django.core.management import utils; print(utils.get_random_secret_key())'"` |
| DJANGO_DEBUG | Bool | Sets django to run in debug mode |
| DJANGO_LANGUAGE_CODE | Str | Language code in RFC 3066 (e.g. "en-us" or "de") |
| DJANGO_TIME_ZONE | Str | e.g. Europe/Berlin |
| DJANGO_PORT | Int | Web port |
| REDIS_HOST | Str | The ip redis runs on |
| REDIS_PORT | Int | The port redis runs on |
| DB_TYPE | Str | Database type. Can be "postgres", "mysql" or "sqlite" |
| DB_HOST | Str | The ip the database runs on |
| DB_PORT | Str | The port the database runs on|
| DB_NAME | Str | Database name |
| DB_USER | Str | Database user |
| DB_PASSWORD | Str | Database password |
| PING_INTERVAL | Int | Time between pings |
| ENABLE_NOTIFICATIONS | Bool | Show notifications in the bottom right corner |

## 📝 Other infos

* The django container needs to run in host network mode to send the wakeonlan command on your local network. Therefore all other containers also need to run in host network mode. I don't like it but there is no way around.
* Firefox 92 and below: The datetime picker for wake events will be available with version 93 and above. [see here](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/datetime-local#browser_compatibility)
