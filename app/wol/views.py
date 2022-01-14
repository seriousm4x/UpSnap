import json
import os
import platform
import subprocess
from io import StringIO

from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django.core import serializers
from django.http import HttpResponseRedirect
from django.http.response import (HttpResponse, HttpResponseServerError,
                                  JsonResponse)
from django.shortcuts import render
from django.utils import dateformat, timezone
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_exempt
from django_wol.settings import VERSION as app_version

from .forms import AddCustomDevice, SettingsForm
from .models import Device, Settings, Websocket


def index(request):
    conf = Settings.objects.first()
    devices = Device.objects.all().order_by(conf.sort_by)
    visitors = Websocket.objects.first().visitors

    context = {
        "devices": devices,
        "visitors": visitors,
        "settings": conf
    }

    return render(request, "wol/index.html", context)


def settings(request):
    conf = Settings.objects.first()
    devices = Device.objects.all().order_by(conf.sort_by)
    visitors = Websocket.objects.first().visitors

    context = {
        "settings": conf,
        "ping_interval": os.getenv("PING_INTERVAL"),
        "devices": devices,
        "platform": platform.uname(),
        "app_version": app_version,
        "debug": os.getenv("DJANGO_DEBUG"),
        "visitors": visitors
    }

    return render(request, "wol/settings.html", context)


def settings_save(request):
    if request.method == "POST":
        form = SettingsForm(request.POST)
        if form.is_valid():
            Settings.objects.update_or_create(
                id=1,
                defaults={
                    "enable_notifications": True if form.cleaned_data["notifications"] == "on" else False,
                    "enable_console_logging": True if form.cleaned_data["console"] == "on" else False,
                    "sort_by": form.cleaned_data["sort"],
                    "scan_address": form.cleaned_data["ip_range"]
                }
            )
    return HttpResponseRedirect("/settings/")


def settings_scan(request):
    data = {
        "devices": []
    }

    conf = Settings.objects.filter(id=1).get()

    if not conf.scan_address:
        return JsonResponse(data=data)

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

            data["devices"].append({
                "name": name,
                "ip": ip,
                "mac": mac
            })

    return JsonResponse(data=data)


def settings_custom_add(request):
    if request.method == "POST":
        form = AddCustomDevice(request.POST)
        if form.is_valid():
            Device.objects.update_or_create(
                mac=form.cleaned_data["custom_add_mac"],
                defaults={
                    "name": form.cleaned_data["custom_add_name"],
                    "ip": form.cleaned_data["custom_add_ip"]
                }
            )
    return HttpResponseRedirect("/settings/")


@method_decorator(csrf_exempt, name="dispatch")
def settings_update(request):
    data = {}
    if request.method == "POST":
        post_body = json.loads(request.body.decode("utf-8"))
        if post_body["name"] == "":
            name = "Unknown"
        else:
            name = post_body["name"]
        Device.objects.update_or_create(
            mac=post_body["mac"],
            defaults={
                "name": name,
                "ip": post_body["ip"]
            }
        )
        data["status"] = 200
    else:
        data["status"] = 500

    channel_layer = get_channel_layer()
    async_to_sync(channel_layer.group_send)(
        "wol", {"type": "send_status", "status": json.dumps({
            "reload": True
        })})

    return JsonResponse(data=data)


@method_decorator(csrf_exempt, name="dispatch")
def settings_del(request):
    data = {}
    if request.method == "POST":
        dev_id = int(request.body.decode("utf-8"))
        Device.objects.get(id=dev_id).delete()
        data["status"] = 200
    else:
        data["status"] = 500

    channel_layer = get_channel_layer()
    async_to_sync(channel_layer.group_send)(
        "wol", {"type": "send_status", "status": json.dumps({
            "reload": True
        })})

    return JsonResponse(data=data)


def settings_export(request):
    if request.method == "GET":
        file_data = {}
        data = serializers.serialize("python", Device.objects.all())
        for dev in data:
            file_data[dev["fields"]["mac"]] = {
                "name": dev["fields"]["name"],
                "ip": dev["fields"]["ip"],
                "netmask": dev["fields"]["netmask"]
            }
        now = dateformat.format(timezone.now(), "Y-m-d_H-i-s")
        response = HttpResponse(
            StringIO(json.dumps(file_data, indent=4)), content_type="application/json")
        response["Content-Disposition"] = f"attachment;filename=upsnap_backup_{now}.json"
        return response
    return HttpResponseRedirect("/settings/")


def settings_import(request):
    if request.method == "POST":
        data = json.loads(request.FILES["file"].read())
        for key, value in data.items():
            Device.objects.update_or_create(
                mac=key,
                defaults={
                    "name": value["name"],
                    "ip": value["ip"],
                    "netmask": value["netmask"],
                }
            )
    return HttpResponseRedirect("/settings/")


def health(request):
    try:
        Settings.objects.all()
        return HttpResponse("Ok")
    except Exception:
        return HttpResponseServerError("db: cannot connect to database.")
