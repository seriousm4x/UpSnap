from django.db import models


class Port(models.Model):
    number = models.PositiveIntegerField()
    name = models.SlugField()

    def __str__(self) -> str:
        return self.name

class Device(models.Model):
    name = models.SlugField(default="Unknown", max_length=100)
    ip = models.GenericIPAddressField()
    mac = models.CharField(max_length=17)
    netmask = models.CharField(max_length=15, default="255.255.255.0", blank=False, null=False)
    scheduled_wake = models.DateTimeField(blank=True, null=True)
    port = models.ManyToManyField(Port, blank=True)

class Websocket(models.Model):
    visitors = models.PositiveSmallIntegerField(blank=False, null=False, default=0)

class Settings(models.Model):
    enable_notifications = models.BooleanField(default=True)
    sort_by = models.SlugField(default="name")
    scan_address = models.GenericIPAddressField(null=True, blank=True)
    interval = models.PositiveSmallIntegerField(null=True, blank=True)
