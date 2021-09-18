from django import forms

from .models import Device


class WakeDeviceForm(forms.ModelForm):
    class Meta:
        model = Device
        fields = []


class AddDeviceForm(forms.ModelForm):
    class Meta:
        model = Device
        fields = []
    device_name = forms.CharField(label="device_name")
    device_ip = forms.CharField(label="device_ip")
    device_mac = forms.CharField(label="device_mac")


class DeleteDeviceForm(forms.ModelForm):
    class Meta:
        model = Device
        fields = []
    device_id = forms.CharField(label="device_id")