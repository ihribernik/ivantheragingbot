import os

from dotenv import load_dotenv
from gtts import gTTS
from twitchio.ext import commands
import pygame
import logging

load_dotenv()


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

    async def event_message(self, message):
        if message.echo:
            return

        await self.handle_commands(message)

    async def reproduce_audio(self, message: str):
        tts = gTTS(text=message, lang=self.lang, slow=False, tld=self.tld)
        tts.save(self.message_name)
        pygame.mixer.music.load(self.message_name)
        pygame.mixer.music.play()
        while pygame.mixer.music.get_busy():
            pygame.time.Clock().tick(10)
        os.remove(self.message_name)

    @commands.command(name="speak")
    async def speak(self, ctx: commands.Context, *, phrase: str):
        message = f"{ctx.author.name} dice: {phrase}"
        self.logger.warning(message)
        await self.reproduce_audio(message)

    async def event_ready(self):
        message = f"Bot has connected to Twitch as {self.nick}"
        self.logger.warning(message)
        await self.reproduce_audio(message)


if __name__ == "__main__":
    bot = ChatReader()
    bot.run()
