import os
from pathlib import Path

from pydantic import SecretStr
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    bot_token: SecretStr
    BASE_DIR: Path = Path(os.path.abspath(__file__))
    nats_url: SecretStr

    class Config:
        env_file = Path(__file__).parent.joinpath(".env")
        env_file_encoding = "utf-8"
        

config = Settings() # type: ignore
