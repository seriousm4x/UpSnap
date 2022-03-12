FROM debian:bullseye-slim as base

FROM base as python-build
ENV PYTHONUNBUFFERED 1
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /python-build
RUN apt-get update &&\
apt-get install -y --no-install-recommends build-essential libffi-dev libssl-dev cargo python3 python3-dev python3-pip python3-venv default-libmysqlclient-dev libpq-dev &&\
python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY app/backend/requirements.txt .
RUN python3 -m pip install --no-cache-dir --upgrade pip wheel &&\
pip install --no-cache-dir -r requirements.txt

FROM node:17-bullseye-slim as node-build
WORKDIR /node-build
COPY app/frontend/package*.json ./
COPY app/frontend/src ./src
COPY app/frontend/public ./public
COPY app/frontend/rollup.config.js ./

FROM base
WORKDIR /app
ENV PYTHONUNBUFFERED 1
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update &&\
apt-get install -y --no-install-recommends default-mysql-client nodejs npm iputils-ping nmap samba-common-bin openssh-client sshpass &&\
apt-get clean &&\
rm -rf /var/lib/{apt,dpkg,cache,log}/
COPY --from=python-build /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY --from=node-build /node-build ./frontend
COPY app/backend ./backend
COPY app/run.sh ./

CMD ["./run.sh"]
