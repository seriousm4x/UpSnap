# Django Wake on LAN

A simple wake on lan app written with Django, Django-Channels (websockets), Celery and Redis.

| Dark                 | Light                 |
| -------------------- | --------------------- |
| ![](assets/dark.jpg) | ![](assets/light.jpg) |

## Run your own instance

You can use the example [docker-compose.yml](docker-compose.yml) file and just run `docker-compose up -d`.

# Available Env Vars

| env var | type | info |
|-|-|-|
| DJANGO_SUPERUSER_USER | Str | Django username for /admin backend |
| DJANGO_SUPERUSER_PASSWORD | Str | Django password for /admin backend |
| DJANGO_SECRET_KEY | Str | Django secrect key. You can create your own by running `docker exec wol_django bash -c "python manage.py shell -c 'from django.core.management import utils; print(utils.get_random_secret_key())'"` |
| DJANGO_DEBUG | Bool | Sets django to run in debug mode |
| DJANGO_LANGUAGE_CODE | Str | Language code in RFC 3066 (e.g. "en-us" or "de) |
| DJANGO_TIME_ZONE | Str | YEP CLOCK |
| DJANGO_PORT | Int | Web port |

# Manage devices

Click on the `Admin` button and login with your django backend user (see [#available-env-vars](#available-env-vars)). Once logged in you can manage your devices there.

# Other infos

* The django container needs to run in host network mode to send the wakeonlan command on your local network. Therefore all other containers also need to run in host network mode. I don't like it but there is no way around.