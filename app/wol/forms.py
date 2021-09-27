from django import forms

from .models import Settings


class SettingsForm(forms.Form):
    class Meta:
        model = Settings
        fields = []
    notifications = forms.CharField(widget=forms.RadioSelect, label="notifications")
    console = forms.CharField(widget=forms.RadioSelect, label="console")
    sort = forms.ChoiceField(choices=[("name", "name"), ("ip", "ip")], widget=forms.RadioSelect, label="sort")
