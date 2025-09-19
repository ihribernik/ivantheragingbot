import logging
import re
from pathlib import Path
from typing import Optional

import pygame
from twitchio.ext import commands

from ivantheraginbot.components.chat import ChatComponent
from ivantheraginbot.components.error import ErrorComponent
from ivantheraginbot.components.help import HelpComponent
from ivantheraginbot.components.sounds import SoundsComponent
from ivantheraginbot.components.speak import SpeakComponent
from ivantheraginbot.config import Settings
from ivantheraginbot.utils import tts


class IvanTheRagingBot(commands.Bot):
    lang: str = "es"
    tld: str = "com.ar"
    virtual_audio_output: str = "CABLE Input (VB-Audio Virtual Cable)"
    audio_output: str = "Digital Audio (S/PDIF) (High Definition Audio Device)"
    url_re: re.Pattern = re.compile(r"https?://(?:www\.)?[^\s/$.?#].[^\s]*")
    settings: Settings
    logger: logging.Logger

    def __init__(
        self, package_location: Path, settings: Optional[Settings] = None
    ) -> None:
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
        await self.add_component(ChatComponent(self))
        await self.add_component(ErrorComponent(self))
        await self.add_component(HelpComponent(self))
        await self.add_component(SoundsComponent(self))
        await self.add_component(SpeakComponent(self))
        return await super().setup_hook()

    async def event_ready(self) -> None:
        username = self.user.name if self.user else "<unknown>"
        message = f"✈️ Bot has connected to Twitch as {username}"
        return await tts(message, self.package_location, self.lang, self.tld)
