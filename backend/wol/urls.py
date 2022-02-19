from django.urls import path, include
from rest_framework import routers
from wol.serializers import SettingsViewSet


router = routers.DefaultRouter()
router.register(r'settings', SettingsViewSet, basename="settings")

urlpatterns = [
    path('api/', include(router.urls)),
]
