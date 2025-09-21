from twitchio import ChatMessage
from twitchio.ext import commands


class ChatComponent(commands.Component):

    def __init__(self, bot: commands.Bot) -> None:
        self.bot = bot

    @commands.Component.listener()
    async def event_message(self, payload: ChatMessage) -> None:
        print(f"[{payload.broadcaster.name}] - {payload.chatter.name}: {payload.text}")
