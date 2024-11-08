import asyncio
import nats


class NATSConnector:
    def __init__(self, url="nats://demo.nats.io:4222", subject="frame_data"):
        self.url = url
        self.subject = subject
        self.nc = None

    async def connect(self):
        # Подключение к NATS
        self.nc = await nats.connect(self.url)
        print(f"Connected to NATS server at {self.url}")

    async def publish_frame_data(self, frame_data):
        # Преобразование данных в JSON и публикация
        payload = json.dumps(frame_data).encode('utf-8')
        await self.nc.publish(self.subject, payload)
        print(f"Published frame_data to '{self.subject}': {frame_data}")

    async def close(self):
        # Закрываем соединение с NATS
        await self.nc.drain()
        print("NATS connection closed.")