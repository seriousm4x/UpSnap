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

# create django secret key
if not os.getenv("DJANGO_SUPERUSER_PASSWORD"):
    os.environ["DJANGO_SECRET_KEY"] = get_random_secret_key()

# reset visitors
[i.delete() for i in Websocket.objects.all()]
Websocket.objects.create(visitors=0)

# notifications
if not os.getenv("PING_INTERVAL"):
    os.environ["PING_INTERVAL"] = 5
Settings.objects.update_or_create(
    id=1,
    defaults={
        "interval": os.getenv("PING_INTERVAL")
    }
)

# register device ping task
schedule, created = IntervalSchedule.objects.get_or_create(
    every=int(os.getenv("PING_INTERVAL")),
    period=IntervalSchedule.SECONDS,
)
if created:
    PeriodicTask.objects.create(
        interval=schedule,
        name="Ping all devices",
        task="wol.tasks.ping_all_devices"
    )
