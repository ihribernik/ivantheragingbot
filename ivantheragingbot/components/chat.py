import logging
from twitchio import ChatMessage
from twitchio.ext import commands

logger = logging.getLogger(__name__)


class ChatComponent(commands.Component):

    def __init__(self, bot: commands.Bot) -> None:
        self.bot = bot

    @commands.Component.listener()
    async def event_message(self, payload: ChatMessage) -> None:
        logger.warning("event_message dispatched")
        msg = (
            f"[{payload.broadcaster.name}] - ",
            f"{payload.chatter.name}: {payload.text}",
        )
        logger.warning(msg)
