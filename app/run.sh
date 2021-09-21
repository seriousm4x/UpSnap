#!/bin/sh

python manage.py makemigrations
python manage.py migrate
python manage.py collectstatic --noinput

# create django super user
if [ -n "$DJANGO_SUPERUSER_USER" ] && [ -n "$DJANGO_SUPERUSER_PASSWORD" ]; then
    python -u manage.py shell -c "from django.contrib.auth.models import User; User.objects.create_superuser('$DJANGO_SUPERUSER_USER', password='$DJANGO_SUPERUSER_PASSWORD') if not User.objects.filter(username='$DJANGO_SUPERUSER_USER').exists() else print('Django user exists')"
fi

# set visitors to 0
    python -u manage.py shell -c "from wol.models import Websocket; [i.delete() for i in Websocket.objects.all()]; Websocket.objects.create(visitors=0)"

celery -A django_wol worker &
celery -A django_wol beat &
gunicorn --bind 0.0.0.0:"$DJANGO_PORT" --workers 4 django_wol.asgi:application -k uvicorn.workers.UvicornWorker

