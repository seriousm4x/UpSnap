import os

from celery import Celery

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "backend.settings")

app = Celery("dango_wol")
app.conf.timezone = os.getenv("DJANGO_TIME_ZONE", "UTC")
app.config_from_object("django.conf:settings", namespace="CELERY")
app.autodiscover_tasks()
