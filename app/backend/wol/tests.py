from django.test import TestCase

# Create your tests here.

from celery.schedules import crontab_parser

print(dir(crontab_parser().parse("0 9 * * 1,5")))
