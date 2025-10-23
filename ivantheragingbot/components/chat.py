import logging

from twitchio import ChatMessage
from twitchio.ext import commands

logger = logging.getLogger(__name__)


class ChatComponent(commands.Component):
    """Component that handles chat messages."""

    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the ChatComponent."""
        super().__init__()
        self.bot = bot

    @commands.Component.listener()
    async def event_message(self, payload: ChatMessage) -> None:
        """Handle chat messages."""
        logger.warning("event_message dispatched")
        msg = (
            f"[{payload.broadcaster.name}] - ",
            f"{payload.chatter.name}: {payload.text}",
        )
        logger.warning(msg)
