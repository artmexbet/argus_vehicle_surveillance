import nats
import json
from telegram_bot.config.config_reader import config


class NATSConnector:
    def __init__(self,url=config.nats_url.get_secret_value(), subject="placeholder") -> None:
        self.url = url
        self.subject = subject
        self.