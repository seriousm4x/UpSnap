import os

from celery import Celery

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "backend.settings")

app = Celery("dango_wol")
app.config_from_object("django.conf:settings", namespace="CELERY")
app.conf.beat_schedule = {
    "ping_devices": {
        "task": "wol.tasks.status",
        "schedule": int(os.getenv("PING_INTERVAL"))
    },
    "scheduled_wakes": {
        "task": "wol.tasks.scheduled_wakes",
        "schedule": 1
    }
}
app.autodiscover_tasks()
