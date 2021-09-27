from django.contrib import admin

from wol.models import Device, Settings


class DeviceAdmin(admin.ModelAdmin):
    list_display = ["name", "ip", "mac"]
    search_fields = ["name", "ip", "mac"]
    list_filter = ["name", "ip", "mac"]


class SettingsAdmin(admin.ModelAdmin):
    list_display = ["enable_notifications",
                    "enable_console_logging", "sort_by"]

    class Meta:
        verbose_name_plural = "Settings"


admin.site.register(Device, DeviceAdmin)
admin.site.register(Settings, SettingsAdmin)
