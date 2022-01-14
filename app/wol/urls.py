from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index"),
    path("settings/", views.settings, name="settings"),
    path("settings/save/", views.settings_save, name="settings_save"),
    path("settings/scan/", views.settings_scan, name="settings_scan"),
    path("settings/update/", views.settings_update, name="settings_update"),
    path("settings/custom_add/", views.settings_custom_add, name="settings_custom_add"),
    path("settings/del/", views.settings_del, name="settings_del"),
    path("settings/export/", views.settings_export, name="settings_export"),
    path("settings/import/", views.settings_import, name="settings_import"),
    path("health/", views.health, name="health")
]
