import logging
import os
import re

import pygame
from gtts import gTTS
from twitchio import Message
from twitchio.ext import commands
from .utils import url_re


class ChatReader(commands.Bot):
    lang = "es"
    MESSAGE_LOCATION = "./message.mp3"
    tld = "com.ar"

    def __init__(self):
        super().__init__(
            token=os.getenv("IRC_TOKEN", "NOT_A_VALID_TOKEN"),
            prefix="!",
            initial_channels=["ivantheragingpython"],
        )
        pygame.mixer.init()
        self.logger = logging.getLogger("ChatReader")

    async def event_ready(self):
        message = f"Bot has connected to Twitch as {self.nick}"
        await self.reproduce_audio(message)

    async def event_message(self, message):
        if message.echo:
            return

        if message.author.name.lower() in self.ignored_users and not os.getenv(
            "READ_AUTHOR_MESSAGE", False
        ):
            return

        parsed_msg = await self.get_context(message)

        if os.getenv("AUTO_READ", None) is not None:
            if parsed_msg.prefix is None:
                message = self.get_parsed_message(message)
                return await self.reproduce_audio(message)

        return await self.handle_commands(message)

    @property
    def ignored_users(self):
        return [self.nick.lower(), "nightbot", "UrMom"]

    def get_parsed_message(self, message: Message) -> str:
        autor = message.author.name or "<unknown>"
        content: str = message.content or ""
        parsed_message = re.sub(url_re, "[Enlace...]", content)
        final_message = f"{autor} dice: {parsed_message}"

        return final_message

    async def event_command_error(
        self,
        context: commands.Context,
        error: Exception,
    ):
        if isinstance(error, commands.CommandNotFound):
            return await context.send(
                "El comando no existe.... !help para ver los commandos disponibles"
            )

        self.logger.error(
            "Error al ejecutar %s: %s",
            context.command.name,
            error,
        )

    async def handle_error(self, ctx, error):
        """Captura y registra errores."""
        self.logger.error("Error en el comando %s:%s", ctx.command, error)
        await ctx.send("Ocurrió un error. Inténtalo de nuevo más tarde.")

    async def reproduce_audio(self, message: str):
        self.logger.warning(message)
        tts = gTTS(text=message, lang=self.lang, slow=False, tld=self.tld)
        tts.save(self.MESSAGE_LOCATION)
        pygame.mixer.music.load(self.MESSAGE_LOCATION)
        pygame.mixer.music.play()
        while pygame.mixer.music.get_busy():
            pygame.time.Clock().tick(10)
        pygame.mixer.music.unload()
        os.remove(self.MESSAGE_LOCATION)

    @property
    def posible_commands(self):
        return ", ".join(self.commands)

    @commands.command(name="speak")
    async def speak(self, ctx: commands.Context, *, phrase: str):
        message = f"{ctx.author.name} dice: {phrase}"
        await self.reproduce_audio(message)

    @commands.command(name="help")
    async def help(self, ctx: commands.Context):
        await ctx.send(f"Comandos disponibles >> {self.posible_commands}")
