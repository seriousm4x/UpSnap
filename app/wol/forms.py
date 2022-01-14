from django import forms

from .models import Settings


class SettingsForm(forms.Form):
    class Meta:
        model = Settings
        fields = []
    notifications = forms.CharField(widget=forms.RadioSelect, label="notifications")
    console = forms.CharField(widget=forms.RadioSelect, label="console")
    sort = forms.ChoiceField(choices=[("name", "name"), ("ip", "ip")], widget=forms.RadioSelect, label="sort")
    ip_range = forms.CharField(label="ip_range", required=False)


class AddCustomDevice(forms.Form):
    custom_add_name = forms.CharField(label="custom_add_name")
    custom_add_ip = forms.CharField(label="custom_add_ip")
    custom_add_mac = forms.CharField(label="custom_add_mac")
