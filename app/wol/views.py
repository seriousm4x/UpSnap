import os
import platform

from django.http import HttpResponseRedirect
from django.shortcuts import render
from django_wol.settings import VERSION as app_version

from .forms import SettingsForm
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
    devices_count = Device.objects.all().count()
    visitors = Websocket.objects.first().visitors


    context = {
        "settings": conf,
        "ping_interval": os.getenv("PING_INTERVAL"),
        "devices_count": devices_count,
        "platform": platform.uname(),
        "app_version": app_version,
        "debug": os.getenv("DJANGO_DEBUG"),
        "visitors": visitors
    }

    return render(request, "wol/settings.html", context)


def save_settings(request):
    if request.method == "POST":
        form = SettingsForm(request.POST)
        if form.is_valid():
            Settings.objects.update_or_create(
                id = 1,
                defaults = {
                    "enable_notifications": True if form.cleaned_data["notifications"] == "on" else False,
                    "enable_console_logging": True if form.cleaned_data["console"] == "on" else False,
                    "sort_by": form.cleaned_data["sort"]
                }
            )
    return HttpResponseRedirect('/settings/')
