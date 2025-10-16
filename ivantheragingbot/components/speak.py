from twitchio.ext import commands

from ivantheragingbot.utils import tts


class SpeakComponent(commands.Component):
    def __init__(self, bot: commands.Bot) -> None:
        self.bot = bot

    @commands.command(name="speak")
    @commands.cooldown(rate=1, per=30, key=commands.BucketType.user)
    async def speak(self, ctx: commands.Context, *, phrase: str) -> None:
        message = f"{ctx.author.name} dice: {phrase}"
        await tts(
            message,
            self.bot.package_location,
            self.bot.lang,
            self.bot.tld,
        )
