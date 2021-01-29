from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("devices/", views.devices, name="devices"),
    path("status/<int:dev_id>", views.status, name="status"),
    path("wake/<int:dev_id>", views.wake, name="wake"),
]
