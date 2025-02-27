import logging
import os

import pygame
from gtts import gTTS
from twitchio.ext import commands


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

        if message.author.name.lower() == self.nick.lower() and not os.getenv(
            "READ_AUTHOR_MESSAGE", False
        ):
            return

        parsed_msg = await self.get_context(message)

        if os.getenv("AUTO_READ", None) is not None:
            if parsed_msg.prefix is None:
                message = f"{message.author.name} dice: {message.content}"
                return await self.reproduce_audio(message)

        return await self.handle_commands(message)

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
