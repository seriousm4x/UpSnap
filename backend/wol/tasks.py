import socket
import subprocess
import threading

from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer
from django_celery_beat.models import PeriodicTask, CrontabSchedule

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
            "ports": [],
            "cron": {
                "enabled": False,
                "value": ""
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
                index = next(i for i, d in enumerate(data["ports"]) if d["number"] == port.number)
                if self.check_port(dev.ip, port.number):
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = True
                else:
                    data["ports"][index]["checked"] = True
                    data["ports"][index]["open"] = False

        # set cron for scheduled wake
        try:
            task = PeriodicTask.objects.filter(name=data["name"], task="wol.tasks.scheduled_wake", crontab_id__isnull=False).get()
            if task:
                cron = CrontabSchedule.objects.get(id=task.crontab_id)
                data["cron"]["enabled"] = task.enabled
                data["cron"]["value"] = " ".join([cron.minute, cron.hour, cron.day_of_week, cron.day_of_month, cron.month_of_year])
        except PeriodicTask.DoesNotExist:
            pass

        async_to_sync(channel_layer.group_send)(
            "wol", {"type": "send_group", "message": {
                "type": "status",
                "message": data
            }})


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

    wake(device.mac, device.ip, device.netmask)
