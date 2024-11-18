FROM python:3.12.4

WORKDIR /app
COPY requirements.txt .
RUN curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
      && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
        sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
        tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
RUN apt-get update && apt-get install ffmpeg libsm6 libxext6 nvidia-container-toolkit  -y
RUN pip install --upgrade pip
RUN pip install -r requirements.txt