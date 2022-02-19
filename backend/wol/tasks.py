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
            "id": dev.id,
            "name": dev.name,
            "ip": dev.ip,
            "mac": dev.mac,
            "netmask": dev.netmask,
            "up": False,
            "ports": []
        }
        for p in Port.objects.all().order_by("number"):
            data["ports"].append({
                "number": p.number,
                "name": p.name,
                "checked": True if p in dev.port.all() else False,
                "open": False
            })

        if self.ping_device(dev.ip):
            data["up"] = True
            for port in dev.port.all():
                index = next(i for i, d in enumerate(data["ports"]) if d["number"] == port.number)
                if self.check_port(dev.ip, port.number):
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = True
                else:
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = False

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
