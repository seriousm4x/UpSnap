FROM python:3-alpine as base

FROM base as builder

ENV PYTHONUNBUFFERED 1

RUN apk update && apk add --no-cache --virtual .build-deps libc-dev make gcc musl-dev python3-dev libffi-dev openssl-dev cargo postgresql-dev
RUN mkdir /install
WORKDIR /install
RUN python -m venv venv
ENV PATH="/install/venv/bin:$PATH"
COPY requirements.txt .
RUN python -m pip install --prefix=/install --no-cache-dir --upgrade pip && \
    pip install --prefix=/install --no-cache-dir -r requirements.txt && \
    apk del .build-deps gcc libc-dev make

FROM base 

COPY --from=builder /install /usr/local
COPY app /app
WORKDIR /app
RUN apk add --no-cache postgresql-dev iputils
RUN addgroup -S user && adduser -S user -G user –no-create-home
RUN chmod -R 755 /app && chown -R user:user /app
USER user

CMD ["./run.sh"]
