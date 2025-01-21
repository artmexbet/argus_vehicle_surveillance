import asyncio
import json
import logging

from aiogram import Bot, Dispatcher, F, types
from aiogram.filters.command import Command
from aiogram.types import ChatMemberUpdated, KeyboardButton, ReplyKeyboardMarkup
from nats_connector import NATSConnector

from telegram_bot.app.messages import (
    HELP_MESSAGE,
    SUBSCRIBE_SUCCESS,
    SUBSCRIBED_MESSAGE,
    UNKNOWN_COMMAND_MESSAGE,
    UNSUBSCRIBE_SUCCESS,
    UNSUBSCRIBED_MESSAGE,
    WELCOME_MESSAGE,
)
from telegram_bot.config.config_reader import config


def create_keyboard(user_telegram_id: int):
    kb_list = [
        [KeyboardButton(text="📋Список команд"), KeyboardButton(text="📖О нас")],
        [KeyboardButton(text="🔔Подписаться"), KeyboardButton(text="❌Отписаться")],
        [KeyboardButton(text="✅Статус подписки")],
    ]
    keyboard = ReplyKeyboardMarkup(keyboard=kb_list, resize_keyboard=True)
    return keyboard


logging.basicConfig(level=logging.INFO)
bot = Bot(token=config.bot_token.get_secret_value())
dp = Dispatcher()


@dp.message(Command("start", "старт"))
@dp.message(F.text == "📖О нас")
async def send_welcome(message: types.Message):
    await message.answer(
        WELCOME_MESSAGE,
        reply_markup=create_keyboard(message.from_user.id),
        parse_mode="HTML",
    )


@dp.message(Command("help"))
@dp.message(F.text == "📋Список команд")
async def help(message: types.Message):
    await message.answer(HELP_MESSAGE, parse_mode="HTML")


@dp.message(Command("подписка"))
@dp.message(F.text == "🔔Подписаться")
async def subscribe(message: types.Message):
    await message.answer(SUBSCRIBE_SUCCESS)


@dp.message(Command("отписка"))
@dp.message(F.text == "❌Отписаться")
async def unsubscribe(message: types.Message):
    await message.answer(UNSUBSCRIBE_SUCCESS)


@dp.message(Command("статус"))
@dp.message(F.text == "✅Статус подписки")
async def status(message: types.Message):
    is_subscribed = True  # Булево значение для статуса пользователя
    if is_subscribed:
        await message.answer(SUBSCRIBED_MESSAGE)
    else:
        await message.answer(UNSUBSCRIBED_MESSAGE)


@dp.message()
async def unknown_command(message: types.Message):
    await message.answer(UNKNOWN_COMMAND_MESSAGE)


async def nats_message_handler(msg):
    try:
        data = msg.data.decode()
        payload = json.loads(data)

        telegram_id = payload.get("telegram_id")
        msg = payload.get("text")

        if not telegram_id:
            print("Ошибка: Telegram ID не найден.")
            return
        text = (
            f"⚠️ Попытка угона автомобиля! ⚠️\n"
            f"Сообщение: {msg}\n"
        )
        await bot.send_message(chat_id=telegram_id, text=text)
        print(f"Уведомление отправлено пользователю: {telegram_id}")
    except Exception as e:
        print(f"Ошибка обработки сообщения: {e}")


async def main():
    nats_connector = NATSConnector(
        nats_url=config.nats_url.get_secret_value(),
        subject=config.nats_subject.get_secret_value(),
        message_handler=nats_message_handler,
    )
    await asyncio.gather(nats_connector.run(), dp.start_polling(bot))
    # await asyncio.gather(dp.start_polling(bot))


if __name__ == "__main__":
    asyncio.run(main())
