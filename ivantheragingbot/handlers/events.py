import logging
import os

from twitchio.ext import commands

from ivantheragingbot.utils.helpers import reproduce_audio


class EventHandler:
    # def __init__(self):
    #     self.logger = logging.getLogger("EventHandler")

    async def event_ready(self):
        message = f"Bot has connected to Twitch as {self.nick}"
        reproduce_audio(message)

    async def event_message(self, message):
        if message.echo:
            return

        if message.author.name.lower() == self.nick.lower():
            return

        parsed_msg = await self.get_context(message)

        if os.getenv("AUTO_READ", None) is not None:
            if parsed_msg.prefix is None:
                message = f"{message.author.name} dice: {message.content}"
                return reproduce_audio(message)

        return await self.handle_commands(message)

    async def event_command_error(self, context: commands.Context, error: Exception):
        if isinstance(error, commands.CommandNotFound):
            return await context.send(
                "El comando no existe.... !help para ver los commandos disponibles"
            )

        self.logger.error("Error al ejecutar %s: %s", context.command.name, error)
