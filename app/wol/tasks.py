import json
import subprocess
import threading

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer

from .models import Device

channel_layer = get_channel_layer()


def ping_device(dev):
    try:
        subprocess.check_output(["ping", "-c", "1", "-W", "1", dev.ip])
        data = {
            "id": dev.id,
            "name": dev.name,
            "ip": dev.ip,
            "mac": dev.mac,
            "up": True
        }
    except subprocess.CalledProcessError:
        data = {
            "id": dev.id,
            "name": dev.name,
            "ip": dev.ip,
            "mac": dev.mac,
            "up": False
        }

    async_to_sync(channel_layer.group_send)(
        "status", {"type": "send_status", "status": json.dumps(data)})


@shared_task
def status():
    devices = Device.objects.all()

    for dev in devices:
        t = threading.Thread(target=ping_device, args=(dev,))
        t.start()
