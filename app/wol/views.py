import ipaddress

import wakeonlan
from django.shortcuts import HttpResponse, get_object_or_404, render

from .forms import WakeDeviceForm
from .models import Device, Websocket


def index(request):
    devices = Device.objects.all().order_by("name")
    visitors = Websocket.objects.first().visitors

    context = {
        "devices": devices,
        "visitors": visitors
    }

    return render(request, "wol/index.html", context)


def wake(request, dev_id):
    dev = get_object_or_404(Device, id=dev_id)
    subnet = ipaddress.ip_network(f"{dev.ip}/{dev.netmask}", strict=False).broadcast_address
    if request.method == "POST":
        form = WakeDeviceForm(request.POST, instance=dev)
        if form.is_valid():
            wakeonlan.send_magic_packet(dev.mac, ip_address=str(subnet))
        return HttpResponse()
