<div align="center" width="100%">
    <img src="app/wol/static/img/favicon.png" width="128" />
</div>

<div align="center" width="100%">
    <h2>UpSnap</h2>
    <p>A simple wake on lan app written with Django, Django-Channels (websockets), Celery, Redis and nmap.</p>
    <a target="_blank" href="https://github.com/seriousm4x/django-wake-on-lan"><img src="https://img.shields.io/github/stars/seriousm4x/django-wake-on-lan" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/pulls/seriousm4x/upsnap" /></a> <a target="_blank" href="https://hub.docker.com/r/seriousm4x/upsnap"><img src="https://img.shields.io/docker/v/seriousm4x/upsnap/latest?label=docker%20image%20ver." /></a> <a target="_blank" href="https://github.com/seriousm4x/django-wake-on-lan"><img src="https://img.shields.io/github/last-commit/seriousm4x/django-wake-on-lan" /></a>
</div>

## ✨ Features

* Dashboard to wake up devices with 1 click
* Set date and time to schedule a wake event
* Open port scan for VNC, RDP and SSH
* Visualization of incoming websocket messages
* Notifications on status changes
* Devices only get pinged when there are 1 or more visitors
* Dark/light mode via preferes-color-scheme
* Settings page to add/delete device and show system infos
* [Docker images](https://hub.docker.com/r/seriousm4x/upsnap) for amd64, arm64, arm/v7

## 📸 Screenshots

| Dark                 | Light                 |
| -------------------- | --------------------- |
| ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/front-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/front-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/settings-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/settings-light.png) |
| ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/schedule-dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/schedule-light.png) |

## 🐳 Run your own instance

You can use the example [docker-compose.yml](docker-compose.yml) file and just run `docker-compose up -d`.

## 🔧 Available Env Vars

| env var | type | info |
|---------|------|------|
| DJANGO_SUPERUSER_USER | Str | Username for /admin backend |
| DJANGO_SUPERUSER_PASSWORD | Str | Password for /admin backend |
| DJANGO_SECRET_KEY | Str | You can create your own by running `docker exec wol_django bash -c "python manage.py shell -c 'from django.core.management import utils; print(utils.get_random_secret_key())'"` |
| DJANGO_DEBUG | Bool | Sets django to run in debug mode |
| DJANGO_LANGUAGE_CODE | Str | Language code in RFC 3066 (e.g. "en-us" or "de") |
| DJANGO_TIME_ZONE | Str | e.g. Europe/Berlin |
| DJANGO_PORT | Int | Web port |
| POSTGRES_USER | Str | Database user |
| POSTGRES_PASSWORD | Str | Database password |
| POSTGRES_DB | Str | Database name |
| PING_INTERVAL | Int | Time between pings |
| ENABLE_NOTIFICATIONS | Bool | Show notifications in the bottom right corner |

## 📝 Other infos

* The django container needs to run in host network mode to send the wakeonlan command on your local network. Therefore all other containers also need to run in host network mode. I don't like it but there is no way around.
