from django.db import models

# Create your models here.

class Device(models.Model):
    name = models.SlugField(max_length=100)
    ip = models.GenericIPAddressField()
    mac = models.CharField(max_length=17)
    netmask = models.CharField(max_length=15, default="255.255.255.0", blank=False, null=False)