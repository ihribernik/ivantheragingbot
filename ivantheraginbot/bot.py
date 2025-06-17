import logging
import os
import re
from pathlib import Path

import pygame
from gtts import gTTS
from twitchio import Message
from twitchio.ext import commands

from .utils import url_re


class ChatReader(commands.Bot):
    lang = "es"
    tld = "com.ar"
    virtual_audio_output = "CABLE Input (VB-Audio Virtual Cable)"
    audio_output = "Digital Audio (S/PDIF) (High Definition Audio Device)"

    def __init__(self, package_location: Path):
        super().__init__(
            token=os.getenv("IRC_TOKEN", "NOT_A_VALID_TOKEN"),
            prefix="!",
            initial_channels=["ivantheragingpython"],
        )
        pygame.mixer.init(devicename=self.audio_output)
        self.logger = logging.getLogger("ChatReader")
        self.package_location = package_location

        self.message_location = self.package_location / "assets/message.mp3"
        self.codec_location = self.package_location / "assets/codec.mp3"
        self.alerta_location = self.package_location / "assets/alerta.mp3"
        self.categoria_location = self.package_location / "assets/categoria.mp3"
        self.can_read = os.getenv("READ_AUTHOR_MESSAGE", None)

    async def event_ready(self):
        message = f"✈️ Bot has connected to Twitch as {self.nick}"
        return await self.tts(message)

    async def event_message(self, message):
        if message.echo:
            return

        if message.author.name.lower() in self.ignored_users and not self.can_read:
            return

        parsed_msg = await self.get_context(message)

        if os.getenv("AUTO_READ", None) is not None:
            if parsed_msg.prefix is None:
                message = self.get_parsed_message(message)
                return await self.tts(message)

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

        if isinstance(error, commands.CommandOnCooldown):
            await context.send(str(error))
            # self.logger.error(error)

        self.logger.error(
            "Error al ejecutar %s: %s",
            context.command.name,
            error,
        )

    async def handle_error(self, ctx, error):
        """Captura y registra errores."""
        self.logger.error("Error en el comando %s:%s", ctx.command, error)
        await ctx.send("Ocurrió un error. Inténtalo de nuevo más tarde.")

    async def reproduce_audio(self, file_location: Path | str):
        pygame.mixer.music.load(file_location)
        pygame.mixer.music.play()
        while pygame.mixer.music.get_busy():
            pygame.time.Clock().tick(10)
        pygame.mixer.music.unload()

    @property
    def posible_commands(self):
        return ", ".join(self.commands)

    async def tts(self, message: str):
        self.logger.warning(message)
        tts = gTTS(text=message, lang=self.lang, slow=False, tld=self.tld)
        tts.save(self.message_location)
        await self.reproduce_audio(self.message_location)
        os.remove(self.message_location)

    @commands.command(name="speak")
    @commands.cooldown(1, 30, commands.Bucket.user)
    async def speak(self, ctx: commands.Context, *, phrase: str):
        message = f"{ctx.author.name} dice: {phrase}"
        self.logger.warning(message)
        await self.tts(message)

    @commands.command(name="help")
    @commands.cooldown(1, 30, commands.Bucket.user)
    async def help(self, ctx: commands.Context):
        await ctx.send(f"🕵️‍♂️ Comandos disponibles >> {self.posible_commands}")

    @commands.command(name="red")
    async def red(self, ctx: commands.Context):
        """comando en homenaje al comoda red de padawanstrainer

        Args:
            ctx (commands.Context): context of the request command
        """
        await self.reproduce_audio(self.codec_location)
        await ctx.send("🛜 Notificacion de red baja enviada...")

    @commands.command(name="alerta")
    @commands.cooldown(1, 30, commands.Bucket.user)
    async def alerat(self, ctx: commands.Context):
        await self.reproduce_audio(self.alerta_location)
        await ctx.send("⚡ ya se alerto al streamer")

    @commands.command(name="categoria")
    @commands.cooldown(1, 30, commands.Bucket.user)
    async def categoria(self, ctx: commands.Context):
        await self.reproduce_audio(self.categoria_location)
        await ctx.send("📷 Se le aviso al streamer que cambie la categoria...")
