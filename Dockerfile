FROM nikolaik/python-nodejs:python3.10-nodejs17-alpine as base

# build python dependencies
FROM base as python-build
WORKDIR /install
ENV PYTHONUNBUFFERED 1
RUN apk update &&\
    apk add musl-dev build-base gcc libffi-dev libressl-dev postgresql-dev mariadb-dev cargo &&\
    rm -rf /var/cache/apk/*
COPY app/backend/requirements.txt .
RUN python -m pip install --no-cache-dir --upgrade pip &&\
    pip install --prefix=/install --no-cache-dir -r requirements.txt

# build svelte app
FROM base as npm-build
WORKDIR /install
COPY app/frontend/package*.json ./
RUN npm install
COPY app/frontend/src ./src
COPY app/frontend/public ./public
COPY app/frontend/rollup.config.js ./
RUN npm run build

# build final image
FROM base
WORKDIR /app
COPY --from=python-build /install /usr/local
COPY app/backend ./backend
COPY --from=npm-build /install ./frontend
COPY app/run.sh ./
RUN apk update &&\
    apk add iputils nmap curl bash mariadb-dev &&\
    rm -rf /var/cache/apk/*

CMD ["./run.sh"]
