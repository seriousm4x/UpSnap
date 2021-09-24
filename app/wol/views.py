from django.shortcuts import render

from wol.models import Device, Websocket


def index(request):
    devices = Device.objects.all().order_by("name")
    visitors = Websocket.objects.first().visitors

    context = {
        "devices": devices,
        "visitors": visitors
    }

    return render(request, "wol/index.html", context)

def settings(request):
    
    context = {
        "test": "ja lol ey"
    }

    return render(request, "wol/settings.html", context)