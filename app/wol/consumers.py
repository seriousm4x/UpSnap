from channels.generic.websocket import AsyncWebsocketConsumer


class WSConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        await self.channel_layer.group_add("status", self.channel_name)
        await self.accept()

    async def disconnect(self, code):
        await self.channel_layer.group_discard("status", self.channel_name)

    async def send_status(self, event):
        text_message = event["text"]

        await self.send(text_message)
