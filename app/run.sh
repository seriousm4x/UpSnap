#!/bin/sh

# init django
cd /app/backend/ || exit
python3 manage.py makemigrations
python3 manage.py migrate
python3 manage.py collectstatic --noinput
python3 manage.py shell < setup.py
python3 -m celery -A backend worker &
python3 -m celery -A backend beat &
python3 -m gunicorn --bind 0.0.0.0:"${BACKEND_PORT}" --workers 4 backend.asgi:application -k uvicorn.workers.UvicornWorker &

# start frontend
cd /app/frontend/ || exit
npm install
npm run build
PORT=$FRONTEND_PORT npm start
