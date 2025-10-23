import logging
import os
from pathlib import Path
from uuid import uuid4

import pygame
from gtts import gTTS

logger = logging.getLogger(__name__)


async def reproduce_audio(file_location: Path | str) -> None:
    pygame.mixer.music.load(file_location)
    pygame.mixer.music.play()
    while pygame.mixer.music.get_busy():
        pygame.time.Clock().tick(10)
    pygame.mixer.music.unload()


async def tts(
    message: str,
    message_location: Path,
    lang: str,
    tld: str,
) -> None:
    logger.warning(message)
    tts = gTTS(
        text=message,
        lang=lang,
        slow=False,
        tld=tld,
    )

    file_location = message_location / str(uuid4())
    tts.save(file_location)
    await reproduce_audio(file_location)
    os.remove(file_location)
