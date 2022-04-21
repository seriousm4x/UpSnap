import socket
import subprocess
import threading

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer
from django.db import connection
from django_celery_beat.models import CrontabSchedule, PeriodicTask

from wol.commands import shutdown, wake
from wol.models import Device, Port, Websocket

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
            "link": dev.link,
            "ports": [],
            "wake": {
                "enabled": False,
                "cron": ""
            },
            "shutdown": {
                "enabled": False,
                "cron": "",
                "command": dev.shutdown_cmd
            }
        }

        # add ports
        for p in Port.objects.all().order_by("number"):
            data["ports"].append({
                "number": p.number,
                "name": p.name,
                "checked": True if p in dev.port.all() else False,
                "open": False
            })

        # set status for device and ports
        if self.ping_device(dev.ip):
            data["up"] = True
            for port in dev.port.all():
                index = next(i for i, d in enumerate(
                    data["ports"]) if d["number"] == port.number)
                if self.check_port(dev.ip, port.number):
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = True
                else:
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = False

        # set cron for wake and shutdown
        for action in ["wake", "shutdown"]:
            try:
                task = PeriodicTask.objects.filter(
                    name=f"{data['name']}-{action}",
                    task=f"wol.tasks.scheduled_{action}", crontab_id__isnull=False).get()
                if task:
                    wake = CrontabSchedule.objects.get(id=task.crontab_id)
                    data[action]["enabled"] = task.enabled
                    data[action]["cron"] = " ".join(
                        [wake.minute, wake.hour, wake.day_of_week, wake.day_of_month, wake.month_of_year])
            except PeriodicTask.DoesNotExist:
                pass

        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_group", "message": {
                "type": "status",
                "message": data
            }})

        connection.close()


@shared_task
def ping_all_devices():
    if Websocket.objects.first().visitors == 0:
        return

    devices = Device.objects.all()

    for dev in devices:
        d = WolDevice()
        t = threading.Thread(target=d.start, args=(dev,))
        t.start()


@shared_task
def scheduled_wake(id):
    try:
        device = Device.objects.get(id=id)
    except Device.DoesNotExist:
        for task in PeriodicTask.objects.filter(args=id):
            task.delete()
        return

    d = WolDevice()
    up = d.ping_device(device.ip)
    if not up:
        wake(device.mac, device.ip, device.netmask)
        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_group", "message": {
                "type": "pending",
                "message": id
            }})


@shared_task
def scheduled_shutdown(id):
    try:
        device = Device.objects.get(id=id)
    except Device.DoesNotExist:
        for task in PeriodicTask.objects.filter(args=id):
            task.delete()
        return

    d = WolDevice()
    up = d.ping_device(device.ip)
    if up:
        shutdown(device.shutdown_cmd)
        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_group", "message": {
                "type": "pending",
                "message": id
            }})
