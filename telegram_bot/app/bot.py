import asyncio
import logging
from aiogram import Bot, Dispatcher, types
from telegram_bot.config.config_reader import config
from aiogram.filters.command import Command
from telegram_bot.app.messages import WELCOME_MESSAGE, HELP_MESSAGE

logging.basicConfig(level=logging.INFO)
bot = Bot(token=config.bot_token.get_secret_value()) 
dp = Dispatcher()


@dp.message(Command("start"))
async def send_welcome(message: types.Message):
    await message.answer(WELCOME_MESSAGE, parse_mode="HTML")


@dp.message(Command("help"))
async def cmd_dice(message: types.Message):
        await message.answer(HELP_MESSAGE, parse_mode="HTML")
    
async def main():
    await dp.start_polling(bot)
    
    
if __name__ == "__main__":
    asyncio.run(main())
