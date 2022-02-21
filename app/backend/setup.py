import os

from django.contrib.auth.models import User

from wol.models import Settings, Websocket
from django_celery_beat.models import PeriodicTask, IntervalSchedule


# create superuser
if not User.objects.filter(username=os.getenv("DJANGO_SUPERUSER_USER")).exists():
    User.objects.create_superuser(os.getenv("DJANGO_SUPERUSER_USER"), password=os.getenv("DJANGO_SUPERUSER_PASSWORD"))
else:
    print('Django user exists')

# reset visitors
[i.delete() for i in Websocket.objects.all()]
Websocket.objects.create(visitors=0)

# notifications
Settings.objects.update_or_create(
    id=1,
    defaults={
        "enable_notifications": os.getenv("ENABLE_NOTIFICATIONS"),
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
