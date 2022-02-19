import socket
import subprocess
import threading

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer
from django.core import serializers
from django.utils import timezone

from wol.models import Device, Port, Websocket
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
                "netmask": dev.netmask,
                "up": False,
                "ports": {}
            }
        }
        for p in Port.objects.all():
            data["device"]["ports"][str(p.number)] = {
                "number": p.number,
                "name": p.name,
                "checked": False,
                "open": False
            }

        if self.ping_device(dev.ip):
            data["device"]["up"] = True
            for port in dev.port.all():
                if self.check_port(dev.ip, port.number):
                    data["device"]["ports"][str(port.number)]["checked"] = True
                    data["device"]["ports"][str(port.number)]["open"] = True
                else:
                    data["device"]["ports"][str(port.number)]["checked"] = True
                    data["device"]["ports"][str(port.number)]["open"] = False

        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_group", "message": {
                "type": "status",
                "message": data
            }})


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
                        "type": "wake_schedule",
                        "message": {
                            "id": dev.id,
                            "name": dev.name
                        }
                    }
                }
            )
