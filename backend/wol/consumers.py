import json

from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django.core import serializers
from django.utils.dateparse import parse_datetime
from django.utils.timezone import make_aware

from wol.models import Device, Port, Websocket
from wol.wake import wake


class WSConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add("wol", self.channel_name)
        await self.accept()
        await self.send(text_data=json.dumps({
            "type": "init",
            "message": await self.get_device_count()
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
                    "visitors": await self.get_visitors()
                }
            }
        )

    async def receive(self, text_data=None):
        received = json.loads(text_data)
        if received["type"] == "wake":
            dev = await self.get_json_from_device_id(received["id"])
            wake(dev["fields"]["mac"], dev["fields"]
                 ["ip"], dev["fields"]["netmask"])

            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "type": "wake",
                        "message": received["id"]
                    }
                }
            )
        elif received["type"] == "add_schedule":
            if not received["datetime"]:
                return
            d = make_aware(parse_datetime(received["datetime"]))
            await self.add_schedule(received["id"], d)
            dev = await self.get_json_from_device_id(received["id"])
            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "add_schedule": {
                            "id": received["id"],
                            "name": dev["fields"]["name"],
                            "datetime": str(d.isoformat())
                        }
                    }
                }
            )
        elif received["type"] == "delete_schedule":
            await self.delete_schedule(received["id"])
            dev = await self.get_json_from_device_id(received["id"])
            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "delete_schedule": {
                            "id": received["id"],
                            "name": dev["fields"]["name"]
                        }
                    }
                }
            )
        elif received["type"] == "get_ports":
            ports = await self.get_ports()
            await self.send(text_data=json.dumps({
                "type": "get_ports",
                "message": ports
            }))
        elif received["type"] == "delete":
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
        elif received["type"] == "update":
            await self.update_device(received["data"])


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
    def get_device_count(self):
        count = Device.objects.count()
        return count

    @database_sync_to_async
    def delete_device(self, id):
        dev = Device.objects.get(id=id)
        dev.delete()

    @database_sync_to_async
    def update_device(self, data):
        obj, _ = Device.objects.update_or_create(
            id=data["id"],
            defaults={
                "name": data["name"],
                "ip": data["ip"],
                "mac": data["mac"],
                "netmask": data["netmask"]
            }
        )
        for _, value in data["ports"].items():
            if value["checked"]:
                p, _ = Port.objects.get_or_create(number=value["number"], name=value["name"])
                obj.port.add(p)
            else:
                p = Port.objects.get(number=value["number"])
                obj.port.remove(p)


    @database_sync_to_async
    def get_json_from_device_id(self, id):
        dev = Device.objects.filter(id=id)
        return serializers.serialize("python", dev)[0]

    @database_sync_to_async
    def add_schedule(self, id, datetime):
        dev = Device.objects.filter(id=id).get()
        dev.scheduled_wake = datetime
        dev.save()

    @database_sync_to_async
    def delete_schedule(self, id):
        dev = Device.objects.filter(id=id).get()
        dev.scheduled_wake = None
        dev.save()

    @database_sync_to_async
    def get_ports(self):
        ports = Port.objects.all().order_by("name")
        return serializers.serialize("python", ports)
