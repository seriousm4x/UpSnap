import json

from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django.core import serializers
from django.utils.dateparse import parse_datetime
from django.utils.timezone import make_aware

from wol.models import Device, Port, Websocket
from wol.wake import wake
from django_celery_beat.models import PeriodicTask, ClockedSchedule, CrontabSchedule


class WSConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add("wol", self.channel_name)
        await self.accept()
        await self.send(text_data=json.dumps({
            "type": "init",
            "message": await self.get_all_devices()
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
            dev = await self.get_device(received["id"])
            wake(dev.mac, dev.ip, dev.netmask)

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
        elif received["type"] == "celery":
            await self.celery_create_scheduled_wake(received["data"])


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
    def get_device(self, id):
        dev = Device.objects.filter(id=id).first()
        return dev


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
