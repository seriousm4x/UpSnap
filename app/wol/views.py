import ipaddress
import json
import subprocess
from platform import system

import wakeonlan
from django.core.serializers import serialize
from django.http import JsonResponse
from django.shortcuts import HttpResponse, get_object_or_404, render

from .forms import WakeDeviceForm
from .models import Device


def index(request):
    devices = Device.objects.all().order_by("name")

    context = {
        "devices": devices
    }

    return render(request, "wol/index.html", context)


def devices(request):
    if request.method == "GET":
        devices = Device.objects.all()
        devices_serialized = json.loads(serialize("json", devices))
        data = {
            "status": 200,
            "devices": devices_serialized
        }
        return JsonResponse(data)


def wake(request, dev_id):
    dev = get_object_or_404(Device, id=dev_id)
    subnet = ipaddress.ip_network(f"{dev.ip}/{dev.netmask}", strict=False).broadcast_address
    if request.method == "POST":
        form = WakeDeviceForm(request.POST, instance=dev)
        if form.is_valid():
            wakeonlan.send_magic_packet(dev.mac, ip_address=str(subnet))
        return HttpResponse()
