from twitchio.ext import commands
from ivantheragingbot.utils import reproduce_audio


class SoundsComponent(commands.Component):

    def __init__(self, bot: commands.Bot) -> None:
        self.bot = bot

    @commands.command(name="red")
    async def red(self, ctx: commands.Context):
        """homenaje al comando red de `padawanstrainer`

        Args:
            ctx (commands.Context): context of the request command
        """
        await reproduce_audio(self.codec_location)
        await ctx.send("🛜 Notificacion de red baja enviada...")

    @commands.command(name="alerta")
    @commands.cooldown(rate=1, per=30, key=commands.BucketType.user)
    async def alerat(self, ctx: commands.Context):
        await reproduce_audio(self.alerta_location)
        await ctx.send("⚡ ya se alerto al streamer")

    @commands.command(name="categoria")
    @commands.cooldown(rate=1, per=30, key=commands.BucketType.user)
    async def categoria(self, ctx: commands.Context):
        await reproduce_audio(self.categoria_location)
        await ctx.send("📷 Se le aviso al streamer que cambie la categoria...")
