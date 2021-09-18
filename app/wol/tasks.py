import json
import subprocess
from platform import system

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer

from .models import Device

channel_layer = get_channel_layer()

@shared_task
def status():
    devices = Device.objects.all()

    data = {
        "status": 200,
        "devices": []
    }

    for dev in devices:
        if system() == "Windows":
            cmd = ["ping", "-n", "1", "-w", "1", dev.ip]
        else:
            cmd = ["ping", "-c", "1", "-W", "1", dev.ip]

        try:
            subprocess.check_output(cmd)
            data["devices"].append({
                "id": dev.id,
                "name": dev.name,
                "ip": dev.ip,
                "mac": dev.mac,
                "up": True
            })
        except subprocess.CalledProcessError:
            data["devices"].append({
                "id": dev.id,
                "name": dev.name,
                "ip": dev.ip,
                "mac": dev.mac,
                "up": False
            })
    
    async_to_sync(channel_layer.group_send)("status", {"type": "send_status", "text": json.dumps(data)})
