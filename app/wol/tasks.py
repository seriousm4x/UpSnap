import json
import socket
import subprocess
import threading

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer
from django.core import serializers
from django.utils import timezone

from wol.models import Device, Websocket
from wol.wake import wake

channel_layer = get_channel_layer()


class WolDevice:
    def ping_device(self, ip):
        try:
            subprocess.check_output(["ping", "-c", "1", "-W", "0.5", ip])
            return True
        except subprocess.CalledProcessError:
            return False

    def check_port(self, ip, port):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(0.5)
        if sock.connect_ex((ip, port)) == 0:
            return True
        return False

    def start(self, dev):
        data = {
            "device": {
                "id": dev.id,
                "name": dev.name,
                "ip": dev.ip,
                "mac": dev.mac,
                "up": False,
                "vnc": False,
                "rdp": False,
                "ssh": False
            }
        }

        if self.ping_device(dev.ip):
            data["device"]["up"] = True
            if self.check_port(dev.ip, 5900):
                data["device"]["vnc"] = True
            if self.check_port(dev.ip, 3389):
                data["device"]["rdp"] = True
            if self.check_port(dev.ip, 22):
                data["device"]["ssh"] = True

        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_status", "status": json.dumps(data)})


@shared_task
def status():
    if Websocket.objects.first().visitors == 0:
        return

    devices = Device.objects.all()

    for dev in devices:
        d = WolDevice()
        t = threading.Thread(target=d.start, args=(dev,))
        t.start()


@shared_task
def scheduled_wakes():
    if Websocket.objects.first().visitors == 0:
        return

    devices = Device.objects.all()

    for dev in devices:
        if dev.scheduled_wake and dev.scheduled_wake <= timezone.now():
            wake(dev.mac, dev.ip, dev.netmask)
            dev.scheduled_wake = None
            dev.save()
            async_to_sync(channel_layer.group_send)(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "wake_schedule": {
                            "id": dev.id,
                            "name": dev.name
                        }
                    }
                }
            )
