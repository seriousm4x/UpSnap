import ipaddress
import json

import wakeonlan
from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncWebsocketConsumer
from django.core import serializers

from .models import Device, Websocket


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
        dev = await self.get_json_from_device_id(text_data)

        subnet = ipaddress.ip_network(
            f"{dev['fields']['ip']}/{dev['fields']['netmask']}", strict=False).broadcast_address
        wakeonlan.send_magic_packet(dev['fields']["mac"], ip_address=str(subnet))

        await self.channel_layer.group_send(
            "wol", {
                "type": "send_group",
                "message": {
                    "wake": dev
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
