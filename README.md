# Django Wake on LAN

A simple wake on lan app written with Django, Django-Channels (websockets), Celery and Redis.

| Dark                 | Light                 |
| -------------------- | --------------------- |
| ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/dark.png) | ![](https://raw.githubusercontent.com/seriousm4x/django-wake-on-lan/master/assets/light.png) |

## Features

* Simple dashboard for waking up devices on local network
* Visual indicators for device up/down and open ports. Pulse effect to visualize incoming websocket messages.
* Schedule wake events
* Notifications on status changes
* Restful pings. Devices only get pinged when there are 1 or more visitors
* Dark/light mode via preferes-color-scheme
* Multithread support for pings
* [Docker images](https://hub.docker.com/r/seriousm4x/django-wol) for amd64, arm64, arm/v7

## Run your own instance

You can use the example [docker-compose.yml](docker-compose.yml) file and just run `docker-compose up -d`.

## Available Env Vars

| env var | type | info |
|---------|------|------|
| DJANGO_SUPERUSER_USER | Str | Django username for /admin backend |
| DJANGO_SUPERUSER_PASSWORD | Str | Django password for /admin backend |
| DJANGO_SECRET_KEY | Str | Django secrect key. You can create your own by running `docker exec wol_django bash -c "python manage.py shell -c 'from django.core.management import utils; print(utils.get_random_secret_key())'"` |
| DJANGO_DEBUG | Bool | Sets django to run in debug mode |
| DJANGO_LANGUAGE_CODE | Str | Language code in RFC 3066 (e.g. "en-us" or "de) |
| DJANGO_TIME_ZONE | Str | YEP CLOCK |
| DJANGO_PORT | Int | Web port |
| POSTGRES_USER | Str | Database user |
| POSTGRES_PASSWORD | Str | Database password |
| POSTGRES_DB | Str | Database name |

## Manage devices

Click on the `Admin` button and login with your django backend user (see [#available-env-vars](#available-env-vars)). Once logged in you can manage your devices there.

## Other infos

* The django container needs to run in host network mode to send the wakeonlan command on your local network. Therefore all other containers also need to run in host network mode. I don't like it but there is no way around.