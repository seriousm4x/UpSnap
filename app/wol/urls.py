from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index"),
    path("settings/", views.settings, name="settings"),
    path("settings/save/", views.settings_save, name="settings_save"),
    path("settings/scan/", views.settings_scan, name="settings_scan"),
    path("settings/scan_add/", views.settings_scan_add, name="settings_scan_add"),
    path("settings/custom_add/", views.settings_custom_add, name="settings_custom_add"),
    path("settings/del/", views.settings_del, name="settings_del")
]
