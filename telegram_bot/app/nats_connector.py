import asyncio
import nats
from telegram_bot.config.config_reader import config


class NATSConnector:
    def __init__(self, message_handler, nats_url=config.nats_url.get_secret_value(), 
                 subject="placeholder", ) -> None:
        self.message_handler = message_handler
        self.nats_url = nats_url
        self.subject = subject
        self.client = None
        
    async def connect(self):
        async def error(e):
            print("Error:", e)
        self.client = await nats.connect(self.nats_url, ping_interval=10, 
                                         verbose=True, error_cb=error)
        print(f"Подключен к NATS: {self.nats_url}")
        
        
    async def subscribe(self):
        if not self.client:
            raise RuntimeError("NATS client не подключён.")
        await self.client.subscribe(self.subject, cb=self.message_handler)
        print(f"Подписан на топик: {self.subject}")
        
        
    async def run(self):
        await self.connect()
        await self.subscribe()
        while True:
            await asyncio.sleep(1)
    