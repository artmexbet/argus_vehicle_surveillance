# Деплой модели
## В докер
1. Собрать образ
```bash
docker build -t python-build -f deploy/python/Dockerfile ./YOLO
```
2. Запустить контейнер
```bash
docker-compose up -d
```
## В локальной среде
1. Установить зависимости
```bash
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu124
pip install -r requirements.txt
```
2. Запустить сервер
```bash
python YOLO/main.py
```