#!/bin/sh

# wait for db and redis
/usr/bin/env bash ./wait-for-it.sh "${POSTGRES_HOST}":"${POSTGRES_PORT}" -t 300 -s
/usr/bin/env bash ./wait-for-it.sh "${REDIS_HOST}":"${REDIS_PORT}" -t 300 -s

python manage.py makemigrations
python manage.py migrate
python manage.py collectstatic --noinput

# create django super user
if [ -n "$DJANGO_SUPERUSER_USER" ] && [ -n "$DJANGO_SUPERUSER_PASSWORD" ]; then
    python -u manage.py shell -c "from django.contrib.auth.models import User; User.objects.create_superuser('$DJANGO_SUPERUSER_USER', password='$DJANGO_SUPERUSER_PASSWORD') if not User.objects.filter(username='$DJANGO_SUPERUSER_USER').exists() else print('Django user exists')"
fi

# set visitors to 0
python -u manage.py shell -c "from wol.models import Websocket; [i.delete() for i in Websocket.objects.all()]; Websocket.objects.create(visitors=0)"

# set notifications
if [ -n "$ENABLE_NOTIFICATIONS" ]; then
    python -u manage.py shell -c "from wol.models import Settings; Settings.objects.update_or_create(id=1, defaults={'enable_notifications': '$ENABLE_NOTIFICATIONS'})"
fi

# set ping interval
if [ -z "$PING_INTERVAL" ]; then
    PING_INTERVAL=5
elif [ "$PING_INTERVAL" -lt 5 ]; then
    echo ""
    echo "Ping interval lower than 5 seconds is not recommended. Please use an interval of 5 seconds or higher. Automatically set to 5 seconds."
    echo ""
    PING_INTERVAL=5
fi

celery -A django_wol worker &
celery -A django_wol beat &
gunicorn --bind 0.0.0.0:"$DJANGO_PORT" --workers 4 django_wol.asgi:application -k uvicorn.workers.UvicornWorker
