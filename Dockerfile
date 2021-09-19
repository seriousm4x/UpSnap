FROM python:3-alpine as base

FROM base as builder

ENV PYTHONUNBUFFERED 1

RUN sed -i 's/dl-cdn.alpinelinux.org/dl-5.alpinelinux.org/g' /etc/apk/repositories && \
    apk update && apk add --no-cache --virtual .build-deps libc-dev make gcc musl-dev python3-dev libffi-dev openssl-dev cargo postgresql-dev && \
    mkdir /install
WORKDIR /install
COPY requirements.txt .
RUN python -m pip install --prefix=/install --no-cache-dir --upgrade pip && \
    pip install --prefix=/install --no-cache-dir -r requirements.txt && \
    apk del .build-deps && \
    rm -rf /usr/local/venv

FROM base

COPY --from=builder /install /usr/local
COPY app /app
WORKDIR /app
RUN apk add --no-cache postgresql-dev iputils && \
    addgroup -S user && adduser -S user -G user –no-create-home && \
    chmod -R 755 /app && chown -R user:user /app
USER user

CMD ["./run.sh"]
