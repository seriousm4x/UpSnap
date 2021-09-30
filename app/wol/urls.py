from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index"),
    path("settings/", views.settings, name="settings"),
    path("settings/save/", views.settings_save, name="settings_save"),
    path("settings/scan/", views.settings_scan, name="settings_scan"),
    path("settings/add/", views.settings_add, name="settings_add"),
    path("settings/del/", views.settings_del, name="settings_del")
]
