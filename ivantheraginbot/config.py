from __future__ import annotations

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=(".env", ".env.prod"),
        env_file_encoding="utf-8",
    )

    READ_AUTHOR_MESSAGE: bool = False
    AUTO_READ: int = 0
    PREFIX: str = "!"
    CLIENT_ID: str = "NOT_A_VALID_CLIENT_ID"
    CLIENT_SECRET: str = "NOT_A_VALID_CLIENT_SECRET"
    BOT_ID: str = "NOT_A_VALID_BOT_ID"
    OWNER_ID: str = "NOT_A_VALID_OWNER_ID"
