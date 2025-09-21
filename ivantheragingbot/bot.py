import logging
import re
import sqlite3
from pathlib import Path
from typing import Optional

import pygame
from asqlite import Pool
from twitchio import MultiSubscribePayload
from twitchio.authentication import UserTokenPayload, ValidateTokenPayload
from twitchio.eventsub import (
    ChatMessageSubscription,
    StreamOnlineSubscription,
    SubscriptionPayload,
)
from twitchio.ext import commands

from ivantheragingbot.components.chat import ChatComponent
from ivantheragingbot.components.error import ErrorComponent
from ivantheragingbot.components.help import HelpComponent
from ivantheragingbot.components.sounds import SoundsComponent
from ivantheragingbot.components.speak import SpeakComponent
from ivantheragingbot.config import Settings
from ivantheragingbot.utils import tts


class IvanTheRagingBot(commands.Bot):
    lang: str = "es"
    tld: str = "com.ar"
    virtual_audio_output: str = "CABLE Input (VB-Audio Virtual Cable)"
    audio_output: str = "Digital Audio (S/PDIF) (High Definition Audio Device)"
    url_re: re.Pattern = re.compile(r"https?://(?:www\.)?[^\s/$.?#].[^\s]*")
    settings: Settings
    logger: logging.Logger

    def __init__(
        self,
        package_location: Path,
        db_pool: Pool,
        settings: Optional[Settings] = None,
    ) -> None:
        self.db_pool = db_pool

        pygame.mixer.init(devicename=self.audio_output)
        self.settings = settings or Settings()
        self.logger = logging.getLogger(__name__)
        self.package_location = package_location

        super().__init__(
            client_id=self.settings.CLIENT_ID,
            client_secret=self.settings.CLIENT_SECRET,
            bot_id=self.settings.BOT_ID,
            owner_id=self.settings.OWNER_ID,
            prefix=self.settings.PREFIX,
        )

    async def setup_hook(self) -> None:
        await self.add_component(HelpComponent(self))
        await self.add_component(ChatComponent(self))
        await self.add_component(ErrorComponent(self))
        await self.add_component(SoundsComponent(self))
        await self.add_component(SpeakComponent(self))

    async def event_oauth_authorized(self, payload: UserTokenPayload) -> None:
        await self.add_token(payload.access_token, payload.refresh_token)

        if not payload.user_id:
            return

        if payload.user_id == self.bot_id:
            # We usually don't want subscribe to events on the bots channel...
            return

        # A list of subscriptions we would like to make to the newly authorized channel...
        subs: list[SubscriptionPayload] = [
            ChatMessageSubscription(
                broadcaster_user_id=payload.user_id, user_id=self.bot_id
            ),
        ]

        resp: MultiSubscribePayload = await self.multi_subscribe(subs)
        if resp.errors:
            self.logger.warning(
                "Failed to subscribe to: %r, for user: %s",
                resp.errors,
                payload.user_id,
            )

    async def add_token(
        self,
        token: str,
        refresh: str,
    ) -> ValidateTokenPayload:
        resp: ValidateTokenPayload = await super().add_token(token, refresh)

        query = """
        INSERT INTO tokens (user_id, token, refresh)
        VALUES (?, ?, ?)
        ON CONFLICT(user_id)
        DO UPDATE SET
            token = excluded.token,
            refresh = excluded.refresh;
        """

        async with self.db_pool.acquire() as connection:
            await connection.execute(query, (resp.user_id, token, refresh))

        self.logger.info(
            "Added token to the database for user: %s",
            resp.user_id,
        )
        return resp

    async def setup_database(
        self,
        db: Pool,
    ) -> tuple[list[tuple[str, str]], list[SubscriptionPayload]]:

        query = """
        CREATE TABLE IF NOT EXISTS tokens(
            user_id TEXT PRIMARY KEY,
            token TEXT NOT NULL,
            refresh TEXT NOT NULL
        )
        """
        async with db.acquire() as connection:
            await connection.execute(query)

            # Fetch any existing tokens...
            rows: list[sqlite3.Row] = await connection.fetchall(
                """SELECT * from tokens"""
            )

            tokens: list[tuple[str, str]] = []
            subs: list[SubscriptionPayload] = []

            for row in rows:
                tokens.append((row["token"], row["refresh"]))

                if row["user_id"] == self.settings.BOT_ID:
                    continue

                subs.extend(
                    [
                        ChatMessageSubscription(
                            broadcaster_user_id=row["user_id"],
                            user_id=self.settings.BOT_ID,
                        )
                    ]
                )

        return tokens, subs

    async def event_ready(self) -> None:
        username = self.user.name if self.user else "<unknown>"
        message = f"✈️ Bot has connected to Twitch as {username}"
        return await tts(message, self.package_location, self.lang, self.tld)
