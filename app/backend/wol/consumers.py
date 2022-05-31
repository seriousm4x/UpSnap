import ipaddress
import json
import os
import subprocess

import pytz
from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django_celery_beat.models import (CrontabSchedule, IntervalSchedule,
                                       PeriodicTask)

from wol.commands import shutdown, wake
from wol.models import Device, Port, Settings, Websocket


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

    async def disconnect(self, _):
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
                        "type": "pending",
                        "message": received["id"]
                    }
                }
            )
        elif received["type"] == "shutdown":
            await self.shutdown_device(received["id"])
            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "type": "pending",
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
            try:
                await self.update_device(received["data"])
                await self.send(text_data=json.dumps({
                    "type": "operationStatus",
                    "message": "Success"
                }))
            except Exception:
                await self.send(text_data=json.dumps({
                    "type": "operationStatus",
                    "message": "Error"
                }))
        elif received["type"] == "update_port":
            await self.update_port(received["data"])
        elif received["type"] == "update_settings":
            await self.update_settings(received["data"])
        elif received["type"] == "scan_network":
            await self.send(text_data=json.dumps({
                "type": "scan_network",
                "message": await self.scan_network()
            }))
        elif received["type"] == "backup":
            await self.send(text_data=json.dumps({
                "type": "backup",
                "message": await self.get_all_devices()
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
    def shutdown_device(self, id):
        dev = Device.objects.filter(id=id).first()
        shutdown(dev.shutdown_cmd)

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
            for p in Port.objects.all().order_by("number"):
                data["ports"].append({
                    "number": p.number,
                    "name": p.name,
                    "checked": False,
                    "open": False
                })
            for action in ["wake", "shutdown"]:
                try:
                    task = PeriodicTask.objects.filter(
                        name=f"{data['name']}-{action}",
                        task=f"wol.tasks.scheduled_{action}", crontab_id__isnull=False).get()
                    if task:
                        wake = CrontabSchedule.objects.get(id=task.crontab_id)
                        data[action]["enabled"] = task.enabled
                        data[action]["cron"] = " ".join(
                            [wake.minute, wake.hour, wake.day_of_month, wake.month_of_year, wake.day_of_week])
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
                "netmask": data["netmask"],
                "link": data["link"] if data.get("link") else "",
                "shutdown_cmd": data["shutdown"]["command"] if data.get("shutdown") else ""
            }
        )
        if data.get("ports"):
            for port in data.get("ports"):
                if port["checked"]:
                    p, _ = Port.objects.get_or_create(
                        number=port["number"], name=port["name"])
                    obj.port.add(p)
                else:
                    p = Port.objects.filter(number=port["number"]).first()
                    if p and p in obj.port.all():
                        obj.port.remove(p)

        for action in ["wake", "shutdown"]:
            if data.get(action):
                if data[action]["enabled"]:
                    cron = data[action]["cron"].strip().split(" ")
                    if not len(cron) == 5:
                        return
                    minute, hour, dom, month, dow = cron
                    schedule, _ = CrontabSchedule.objects.get_or_create(
                        minute=minute,
                        hour=hour,
                        day_of_week=dow,
                        day_of_month=dom,
                        month_of_year=month,
                        timezone = pytz.timezone(os.getenv("DJANGO_TIME_ZONE", "UTC"))
                    )
                    PeriodicTask.objects.update_or_create(
                        name=f"{data['name']}-{action}",
                        defaults={
                            "crontab": schedule,
                            "task": f"wol.tasks.scheduled_{action}",
                            "args": json.dumps([data["id"]]),
                            "enabled": True
                        }
                    )
                else:
                    for task in PeriodicTask.objects.filter(name=f"{data['name']}-{action}", task=f"wol.tasks.scheduled_{action}"):
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
            "discovery": conf.scan_address,
            "interval": conf.interval,
            "scan_network": [],
            "notifications": conf.notifications
        }
        return data

    @database_sync_to_async
    def update_settings(self, data):
        if data["interval"] > 5:
            data["interval"] = 5
        Settings.objects.update_or_create(
            id=1,
            defaults={
                "scan_address": data["discovery"],
                "interval": data["interval"],
                "notifications": data["notifications"]
            }
        )
        schedule, _ = IntervalSchedule.objects.get_or_create(
            every=data["interval"],
            period=IntervalSchedule.SECONDS,
        )
        PeriodicTask.objects.update_or_create(
            task="wol.tasks.ping_all_devices",
            defaults={
                "interval": schedule
            }
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
                    "mac": mac,
                    "link": "",
                    "ports": [],
                    "wake": {
                        "enabled": False,
                        "cron": ""
                    },
                    "shutdown": {
                        "enabled": False,
                        "cron": "",
                        "command": ""
                    }
                })
        return data
