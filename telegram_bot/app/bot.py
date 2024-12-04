from aiogram import Bot
from telegram_bot.config.config_reader import config


bot = Bot(token=config.bot_token.get_secret_value()) 
