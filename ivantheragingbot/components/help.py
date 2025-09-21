from twitchio.ext import commands


class HelpComponent(commands.Component):

    def __init__(self, bot: commands.Bot) -> None:
        self.bot = bot

    @commands.command(name="help")
    async def help(self, ctx: commands.Context):
        print("calling the help command")
        await ctx.send("🕵️‍♂️ Comandos disponibles >>")
