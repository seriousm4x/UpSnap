from rest_framework import serializers, viewsets
from wol.models import Settings, Device, Port

class SettingsSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Settings
        fields = "__all__"

class SettingsViewSet(viewsets.ModelViewSet):
    def get_queryset(self):
        queryset = Settings.objects.first()
        return queryset
    serializer_class = SettingsSerializer

class PortSerializer(serializers.ModelSerializer):
    class Meta:
        model = Port
        fields = ["number", "name"]

class DeviceSerializer(serializers.ModelSerializer):
    port = PortSerializer(many=True, read_only=True)

    class Meta:
        model = Device
        fields = ["name", "ip", "mac", "netmask", "scheduled_wake", "port"]
