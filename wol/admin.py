from django.contrib import admin

from .models import Device


class DeviceAdmin(admin.ModelAdmin):
    list_display = ["name", "ip", "mac"]
    search_fields = ["name", "ip", "mac"]
    list_filter = ["name", "ip", "mac"]

admin.site.register(Device, DeviceAdmin)
