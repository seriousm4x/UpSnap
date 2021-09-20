FROM python:3-slim-bullseye as base

FROM base as builder

ENV PYTHONUNBUFFERED 1

RUN apt-get update && apt-get -y install build-essential libssl-dev libffi-dev python3-dev cargo libpq-dev && \
    mkdir /install
WORKDIR /install
COPY requirements.txt .
RUN python -m pip install --prefix=/install --no-cache-dir --upgrade pip && \
    pip install --prefix=/install --no-cache-dir -r requirements.txt && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

FROM base

COPY --from=builder /install /usr/local
COPY app /app
WORKDIR /app
RUN apt-get -y install iputils-ping && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd user && useradd -M user -g user && \
    chmod -R 755 /app && chown -R user:user /app
USER user

CMD ["./run.sh"]
