from array import array
import ipaddress
import json
import subprocess

from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django.utils.dateparse import parse_datetime
from django.utils.timezone import make_aware
from django_celery_beat.models import (ClockedSchedule, CrontabSchedule,
                                       PeriodicTask)

from wol.models import Device, Port, Settings, Websocket
from wol.wake import wake


class WSConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add("wol", self.channel_name)
        await self.accept()
        await self.send(text_data=json.dumps({
            "type": "init",
            "message": {
                "devices": await self.get_all_devices(),
                "settings": await self.get_settings()
            }
        }))
        await self.add_visitor()
        await self.channel_layer.group_send(
            "wol", {
                "type": "send_group",
                "message": {
                    "type": "visitor",
                    "message": await self.get_visitors()
                }
            }
        )

    async def disconnect(self, code):
        await self.channel_layer.group_discard("wol", self.channel_name)
        await self.remove_visitor()
        await self.channel_layer.group_send(
            "wol", {
                "type": "send_group",
                "message": {
                    "type": "visitor",
                    "message": await self.get_visitors()
                }
            }
        )

    async def receive(self, text_data=None):
        received = json.loads(text_data)
        if received["type"] == "wake":
            await self.wake_device(received["id"])
            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "type": "wake",
                        "message": received["id"]
                    }
                }
            )
        elif received["type"] == "delete_device":
            await self.delete_device(received["id"])
            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "type": "delete",
                        "message":  received["id"]
                    }
                }
            )
        elif received["type"] == "update_device":
            await self.update_device(received["data"])
        elif received["type"] == "update_port":
            await self.update_port(received["data"])
        elif received["type"] == "update_settings":
            await self.update_settings(received["data"])
        elif received["type"] == "celery":
            await self.celery_create_scheduled_wake(received["data"])
        elif received["type"] == "scan_network":
            await self.send(text_data=json.dumps({
                "type": "scan_network",
                "message": await self.scan_network()
            }))

    async def send_group(self, event):
        await self.send(json.dumps(event["message"]))

    @database_sync_to_async
    def add_visitor(self):
        visitor_count = Websocket.objects.first()
        visitor_count.visitors += 1
        visitor_count.save()

    @database_sync_to_async
    def remove_visitor(self):
        visitor_count = Websocket.objects.first()
        visitor_count.visitors -= 1
        visitor_count.save()

    @database_sync_to_async
    def get_visitors(self):
        return Websocket.objects.first().visitors

    @database_sync_to_async
    def wake_device(self, id):
        dev = Device.objects.filter(id=id).first()
        wake(dev.mac, dev.ip, dev.netmask)

    @database_sync_to_async
    def get_all_devices(self):
        devices = Device.objects.all()
        d = []
        for dev in devices:
            data = {
                "id": dev.id,
                "name": dev.name,
                "ip": dev.ip,
                "mac": dev.mac,
                "netmask": dev.netmask,
                "ports": [],
                "cron": {
                    "enabled": False,
                    "value": ""
                }
            }
            for p in Port.objects.all().order_by("number"):
                data["ports"].append({
                    "number": p.number,
                    "name": p.name,
                    "checked": False,
                    "open": False
                })
            try:
                task = PeriodicTask.objects.filter(name=data["name"], task="wol.tasks.scheduled_wake", crontab_id__isnull=False).get()
                if task:
                    cron = CrontabSchedule.objects.get(id=task.crontab_id)
                    data["cron"]["enabled"] = task.enabled
                    data["cron"]["value"] = " ".join([cron.minute, cron.hour, cron.day_of_week, cron.day_of_month, cron.month_of_year])
            except PeriodicTask.DoesNotExist:
                pass
            d.append(data)
        return d

    @database_sync_to_async
    def delete_device(self, id):
        dev = Device.objects.get(id=id)
        dev.delete()

    @database_sync_to_async
    def update_device(self, data):
        obj, _ = Device.objects.update_or_create(
            mac=data["mac"],
            defaults={
                "name": data["name"],
                "ip": data["ip"],
                "netmask": data["netmask"]
            }
        )
        if data.get("ports"):
            for port in data.get("ports"):
                if port["checked"]:
                    p, _ = Port.objects.get_or_create(number=port["number"], name=port["name"])
                    obj.port.add(p)
                else:
                    p = Port.objects.get(number=port["number"])
                    obj.port.remove(p)

        if data["cron"]["enabled"]:
            cron_value = data["cron"]["value"].strip().split(" ")
            if not len(cron_value) == 5:
                return
            minute, hour, dom, month, dow = cron_value
            schedule, _ = CrontabSchedule.objects.get_or_create(
                minute=minute,
                hour=hour,
                day_of_week=dow,
                day_of_month=dom,
                month_of_year=month
            )
            PeriodicTask.objects.update_or_create(
                name=data["name"],
                defaults={
                    "crontab": schedule,
                    "task": "wol.tasks.scheduled_wake",
                    "args": json.dumps([data["id"]]),
                    "enabled": True
                }
            )
        else:
            for task in PeriodicTask.objects.filter(name=data["name"], task="wol.tasks.scheduled_wake"):
                task.enabled = False
                task.save()

    @database_sync_to_async
    def update_port(self, data):
        if data.get("name"):
            Port.objects.update_or_create(
                number=data["number"],
                defaults={
                    "name": data["name"]
                }
            )
        else:
            Port.objects.filter(number=data["number"]).delete()

    @database_sync_to_async
    def get_settings(self):
        conf = Settings.objects.get(id=1)
        data = {
            "notifications": conf.enable_notifications,
            "discovery": conf.scan_address,
            "interval": conf.interval,
            "scan_network": []
        }
        return data

    @database_sync_to_async
    def update_settings(self, data):
        Settings.objects.update_or_create(
            id=1,
            defaults={
                "enable_notifications": data["notifications"],
                "scan_address": data["discovery"],
                "interval": data["interval"]
            }
        )

    @database_sync_to_async
    def celery_create_scheduled_wake(self, data):
        aware_time = make_aware(parse_datetime(data["time"]))
        schedule, _ = ClockedSchedule.objects.get_or_create(
            clocked_time=aware_time
        )
        PeriodicTask.objects.get_or_create(
            clocked=schedule,
            name=data["name"],
            task="wol.tasks.scheduled_wake",
            args=json.dumps([data["id"]])
        )

    @database_sync_to_async
    def scan_network(self):
        data = []
        conf = Settings.objects.get(id=1)
        netmask = str(ipaddress.ip_network(conf.scan_address).netmask)

        if not conf.scan_address:
            return

        p = subprocess.Popen(
            ["nmap", "-sP", conf.scan_address], stdout=subprocess.PIPE)
        out = p.communicate()[0].decode("utf-8")
        ip_line = "Nmap scan report for"
        mac_line = "MAC Address:"

        for line in out.splitlines():
            if line.startswith(ip_line):
                line_splitted = line.split()
                if len(line_splitted) == 6:
                    name = line_splitted[4]
                    ip = line_splitted[5]
                    ip = ip.replace("(", "")
                    ip = ip.replace(")", "")
                else:
                    name = "Unknown"
                    ip = line_splitted[4]
            elif line.startswith(mac_line):
                line_splitted = line.split()
                mac = line_splitted[2]
                data.append({
                    "name": name,
                    "ip": ip,
                    "netmask": netmask,
                    "mac": mac
                })
        return data
