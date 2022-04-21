import os

from django.contrib.auth.models import User
from django_celery_beat.models import IntervalSchedule, PeriodicTask, PeriodicTasks

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
    if int(os.environ.get("PING_INTERVAL")) < 5:
        print("Ping interval lower than 5 seconds is not recommended. Please use an interval of 5 seconds or higher. Automatically set to 5 seconds.")
        ping_interval = 5
    else:
        ping_interval = int(os.environ.get("PING_INTERVAL"))
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

# reset last run to fix time zone changes
# https://django-celery-beat.readthedocs.io/en/latest/#important-warning-about-time-zones
PeriodicTask.objects.update(last_run_at=None)
for task in PeriodicTask.objects.all():
    PeriodicTasks.changed(task)
