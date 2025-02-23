import os

import pygame
from gtts import gTTS
import logging

logger = logging.getLogger(__name__)
MESSAGE_LOCATION = os.getenv("MESSAGE_LOCATION", "./message.mp3")


def reproduce_audio(message: str):
    logger.warning(message)
    tts = gTTS(text=message, lang=self.lang, slow=False, tld=self.tld)
    tts.save(MESSAGE_LOCATION)
    pygame.mixer.music.load(MESSAGE_LOCATION)
    pygame.mixer.music.play()
    while pygame.mixer.music.get_busy():
        pygame.time.Clock().tick(10)
    pygame.mixer.music.unload()
    os.remove(MESSAGE_LOCATION)


def get_commands_bots(cls):
    return cls.posible_commands
