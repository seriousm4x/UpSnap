# Generated by Django 3.2.12 on 2022-02-20 19:53

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('wol', '0002_remove_settings_enable_console_logging'),
    ]

    operations = [
        migrations.AddField(
            model_name='settings',
            name='interval',
            field=models.PositiveSmallIntegerField(blank=True, null=True),
        ),
    ]