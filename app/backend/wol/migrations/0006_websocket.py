# Generated by Django 3.2.7 on 2021-09-22 20:02

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('wol', '0005_device_netmask'),
    ]

    operations = [
        migrations.CreateModel(
            name='Websocket',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('visitors', models.PositiveSmallIntegerField(default=0)),
            ],
        ),
    ]