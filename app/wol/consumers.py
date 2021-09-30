import json
from datetime import datetime, tzinfo
from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django.core import serializers
from django.utils.dateparse import parse_datetime
from django.utils.timezone import make_aware

from wol.models import Device, Websocket
from wol.wake import wake


class WSConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add("wol", self.channel_name)
        await self.accept()
        await self.add_visitor()
        await self.channel_layer.group_send(
            "wol", {
                "type": "send_group",
                "message": {
                    "visitors": await self.get_visitors()
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

    async def receive(self, text_data=None, bytes_data=None):
        received = json.loads(text_data)
        if received["message"] == "wake":
            dev = await self.get_json_from_device_id(received["id"])
            wake(dev["fields"]["mac"], dev["fields"]
                 ["ip"], dev["fields"]["netmask"])

            await self.channel_layer.group_send(
                "wol", {
                    "type": "send_group",
                    "message": {
                        "wake": {
                            "id": dev["pk"],
                            "name": dev["fields"]["name"]
                        }
                    }
                }
            )
        elif received["message"] == "add_schedule":
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
        elif received["message"] == "delete_schedule":
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


    async def send_status(self, event):
        await self.send(event["status"])

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
