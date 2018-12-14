from app import app
from celery.bin import purge


if __name__ == '__main__':
    purge.purge()  # Purge Celery
    app.app.run(host="0.0.0.0", port=80)  # Initialize app
