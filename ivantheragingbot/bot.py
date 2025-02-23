import logging
import os

import pygame
from twitchio.ext import commands

from ivantheragingbot.handlers import (CommandsHandler, ErrorHandler,
                                       EventHandler)


class ChatReader(commands.Bot):
    lang = "es"
    message_name = "./message.mp3"
    tld = "com.ar"

    def __init__(self):
        super().__init__(
            token=os.getenv("IRC_TOKEN", "NOT_A_VALID_TOKEN"),
            prefix="!",
            initial_channels=["ivantheragingpython"],
        )
        pygame.mixer.init()
        self.logger = logging.getLogger("ChatReader")

        # self.command_handler = CommandsHandler()
        # self.event_handler = EventHandler()
        # self.error_handler = ErrorHandler()

    @property
    def posible_commands(self):
        return ", ".join(self.commands)
