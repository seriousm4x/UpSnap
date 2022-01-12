FROM python:3.10-slim-bullseye as base

FROM base as builder

ENV PYTHONUNBUFFERED 1

RUN apt-get update && apt-get -y install build-essential libssl-dev libffi-dev python3-dev cargo libpq-dev && \
    mkdir /install
WORKDIR /install
COPY requirements.txt .
RUN python -m pip install --no-cache-dir --upgrade pip && \
    pip install --prefix=/install --no-cache-dir -r requirements.txt && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

FROM base

COPY --from=builder /install /usr/local
COPY app /app
WORKDIR /app
RUN apt-get update && \
    apt-get -y install iputils-ping nmap curl && \
    rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:$DJANGO_PORT/health/" || exit 1

CMD ["./run.sh"]
