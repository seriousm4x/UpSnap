from django.contrib import admin

from wol.models import Device, Port, Settings


class DeviceAdmin(admin.ModelAdmin):
    list_display = ["name", "ip", "mac", "netmask", "ports"]
    search_fields = ["name", "ip", "mac"]
    list_filter = ["name", "netmask"]

    def ports(self, obj):
        return ", ".join([p.name for p in obj.port.all()])

class PortAdmin(admin.ModelAdmin):
    list_display = ["number", "name"]
    search_fields = ["number", "name"]
    list_filter = ["number", "name"]

class SettingsAdmin(admin.ModelAdmin):
    list_display = ["enable_notifications", "sort_by", "scan_address", "interval"]


admin.site.register(Device, DeviceAdmin)
admin.site.register(Port, PortAdmin)
admin.site.register(Settings, SettingsAdmin)
