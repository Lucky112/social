FROM python:3.9-slim

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir psycopg2-binary

COPY insert_data.py /app/insert_data.py
COPY data.csv /app/data.csv

WORKDIR /app

# provided in docker-compose
ENV DB_HOST=""
ENV DB_PORT=""
ENV DB_NAME=""
ENV DB_USER=""
ENV DB_PASSWORD_FILE=""

CMD ["python", "insert_data.py"]
