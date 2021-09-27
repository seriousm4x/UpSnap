from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index"),
    path("settings/", views.settings, name="settings"),
    path("save_settings/", views.save_settings, name="save_settings")
]
