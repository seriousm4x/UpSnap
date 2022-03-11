import os

from django.contrib.auth.models import User
from django.core.management.utils import get_random_secret_key
from django_celery_beat.models import IntervalSchedule, PeriodicTask

from wol.models import Settings, Websocket

# create superuser
if os.getenv("DJANGO_SUPERUSER_USER") and os.getenv("DJANGO_SUPERUSER_PASSWORD"):
    if not User.objects.filter(username=os.getenv("DJANGO_SUPERUSER_USER")).exists():
        User.objects.create_superuser(os.getenv("DJANGO_SUPERUSER_USER"), password=os.getenv("DJANGO_SUPERUSER_PASSWORD"))
    else:
        print('Django user exists')

# reset visitors
[i.delete() for i in Websocket.objects.all()]
Websocket.objects.create(visitors=0)

# ping interval
if os.environ.get("PING_INTERVAL"):
    ping_interval = os.environ.get("PING_INTERVAL")
else:
    ping_interval = 5
Settings.objects.update_or_create(
    id=1,
    defaults={
        "interval": ping_interval
    }
)

# register device ping task
schedule, created = IntervalSchedule.objects.get_or_create(
    every=ping_interval,
    period=IntervalSchedule.SECONDS,
)
if created:
    PeriodicTask.objects.create(
        interval=schedule,
        name="Ping all devices",
        task="wol.tasks.ping_all_devices"
    )
