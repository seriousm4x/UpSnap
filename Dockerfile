FROM python:3.10-alpine as base

FROM base as builder

ENV PYTHONUNBUFFERED 1

RUN apk update &&\
    apk add python3-dev musl-dev build-base gcc libffi-dev libressl-dev postgresql-dev cargo &&\
    rm -rf /var/cache/apk/* &&\
    mkdir /install
WORKDIR /install
COPY requirements.txt .
RUN python -m pip install --no-cache-dir --upgrade pip &&\
    pip install --prefix=/install --no-cache-dir -r requirements.txt

FROM base

COPY --from=builder /install /usr/local
COPY app /app
WORKDIR /app
RUN apk update &&\
    apk add iputils nmap curl bash &&\
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:$DJANGO_PORT/health/" || exit 1

CMD ["./run.sh"]
