# build svelte app
FROM node:17-alpine as npm-build
WORKDIR /install
COPY app/frontend/package*.json ./
RUN npm install
COPY app/frontend/src ./src
COPY app/frontend/public ./public
COPY app/frontend/rollup.config.js ./
RUN npm run build

# build python dependencies
FROM python:3.10-alpine
ENV PYTHONUNBUFFERED 1
WORKDIR /app
RUN apk update &&\
    apk add musl-dev build-base gcc libffi-dev libressl-dev postgresql-dev mariadb-dev nodejs npm iputils nmap curl bash &&\
    rm -rf /var/cache/apk/*
COPY app/backend/requirements.txt .
RUN python -m pip install --no-cache-dir --upgrade pip &&\
    pip install --no-cache-dir -r requirements.txt &&\
    apk del build-base gcc libffi-dev libressl-dev postgresql-dev
COPY app/backend ./backend
COPY --from=npm-build /install ./frontend
COPY app/run.sh ./

CMD ["./run.sh"]
