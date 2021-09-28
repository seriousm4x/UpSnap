from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index"),
    path("settings/", views.settings, name="settings"),
    path("save-settings/", views.save_settings, name="save_settings"),
    path("scan/", views.scan, name="scan"),
    path("add-device/", views.add_device, name="add_device")
]
