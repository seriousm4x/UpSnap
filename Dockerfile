FROM python:3-alpine

ENV PYTHONUNBUFFERED 1

WORKDIR /opt/app
COPY . .

RUN python -m pip install --upgrade pip && pip install -r requirements.txt
