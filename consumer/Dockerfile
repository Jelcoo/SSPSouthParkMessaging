FROM python:3.13-alpine
WORKDIR /app/consumer

ENV PYTHONUNBUFFERED=1

COPY consumer/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY consumer/ ./

CMD ["python", "main.py"]
