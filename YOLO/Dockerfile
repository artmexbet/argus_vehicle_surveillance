FROM python-build
WORKDIR /app
COPY . .
RUN mkdir -p YOLO  # Удостовериться, что работает так
COPY yolo11l.pt YOLO
CMD ["python", "main.py"]