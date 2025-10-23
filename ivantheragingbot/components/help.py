from __future__ import annotations

from twitchio.ext import commands

from ivantheragingbot.types import BotT


class HelpComponent(commands.Component):

    def __init__(self, bot: BotT) -> None:
        super().__init__()
        self.bot: BotT = bot

    @commands.command(name="help")
    async def help(self, ctx: commands.Context):
        self.bot.logger.warning("calling the help command")
        await ctx.send("🕵️‍♂️ Comandos disponibles >>")
