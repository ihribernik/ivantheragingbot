from __future__ import annotations

import asyncio
from typing import Optional

from twitchio.ext import commands

from ivantheragingbot.types import BotT
from ivantheragingbot.utils import tts


class SpeakComponent(commands.Component):
    """Component that handles text-to-speech commands."""

    def __init__(self, bot: BotT) -> None:
        super().__init__()
        self.bot: BotT = bot

    @commands.command(name="speak")
    @commands.cooldown(rate=1, per=30, key=commands.BucketType.user)
    async def speak(
        self,
        ctx: commands.Context,
        *,
        phrase: Optional[str] = None,
    ) -> None:
        """Speak a user-provided phrase using TTS."""
        if not phrase:
            await ctx.send(
                f"{ctx.author.name}, you need to provide something to say!",
            )
            return

        message = f"{ctx.author.name} dice: {phrase}"

        try:
            loop = asyncio.get_running_loop()
            result = await loop.run_in_executor(
                None,
                tts,
                message,
                self.bot.package_location,
                self.bot.lang,
                self.bot.tld,
            )

            if result:
                await ctx.send("🗣️ Message spoken.")
            else:
                await ctx.send("⚠️ Error generating TTS.")
        except Exception as e:
            await ctx.send("⚠️ Error generating TTS.")
            self.bot.logger.exception("Error in !speak command", exc_info=e)

    async def on_error(self, ctx: commands.Context, error: Exception) -> None:
        """Local component-specific error handler."""
        await ctx.send(f"⚠️ SpeakComponent error: {error}")
