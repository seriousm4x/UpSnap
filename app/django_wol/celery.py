import os
from celery import Celery

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "django_wol.settings")

app = Celery("dango_wol")

app.config_from_object("django.conf:settings", namespace="CELERY")

app.conf.beat_schedule = {
    "ping_devices_5s": {
        "task": "wol.tasks.status",
        "schedule": 5
    },
    "scheduled_wakes_1s": {
        "task": "wol.tasks.scheduled_wakes",
        "schedule": 1
    }
}

app.autodiscover_tasks()
