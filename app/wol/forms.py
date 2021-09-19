from django import forms

from .models import Device


class WakeDeviceForm(forms.ModelForm):
    class Meta:
        model = Device
        fields = []
