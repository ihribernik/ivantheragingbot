from __future__ import annotations

import logging

from twitchio.ext import commands

from ivantheragingbot.types import BotT

logger = logging.getLogger(__name__)


class ErrorComponent(commands.Component):

    def __init__(self, bot: BotT) -> None:
        super().__init__()
        self.original = bot.event_command_error
        bot.event_command_error = self.event_command_error
        self.bot = bot

    async def component_teardown(self) -> None:
        self.bot.event_command_error = self.original

    async def event_command_error(
        self,
        payload: commands.CommandErrorPayload,
    ) -> None:
        ctx = payload.context
        command = ctx.command
        error = payload.exception

        if command and command.has_error and ctx.error_dispatched:
            return

        if isinstance(error, commands.CommandNotFound):
            return

        if isinstance(error, commands.GuardFailure):
            await ctx.send(
                str(ctx.chatter) +
                " you are not allowed to use this command."
            )

        msg = f'Ignoring exception in command "{ctx.command}":\n'
        logger.error(msg, exc_info=error)
