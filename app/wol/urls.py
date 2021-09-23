from django.urls import path

from wol import views

urlpatterns = [
    path("", views.index, name="index")
]
