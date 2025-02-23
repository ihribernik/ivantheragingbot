from twitchio.ext import commands

from ivantheragingbot.utils.helpers import reproduce_audio, get_commands_bots


class CommandsHandler:
    @commands.command(name="speak")
    async def speak(self, ctx: commands.Context, *, phrase: str):
        message = f"{ctx.author.name} dice: {phrase}"
        reproduce_audio(message)

    @commands.command(name="help")
    async def help(self, ctx: commands.Context):
        await ctx.send(f"Comandos disponibles >> {get_commands_bots(self)}")
