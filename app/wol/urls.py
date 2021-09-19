from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("wake/<int:dev_id>", views.wake, name="wake"),
]
