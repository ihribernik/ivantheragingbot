import asyncio
import logging
from pathlib import Path

import asqlite
import twitchio

from ivantheragingbot.bot import IvanTheRagingBot
from ivantheragingbot.config import Settings

LOGGER: logging.Logger = logging.getLogger("ivantheraginbot")


def main() -> None:
    package_location = Path.cwd()
    twitchio.utils.setup_logging(level=logging.WARNING)
    settings = Settings()

    async def runner() -> None:
        async with asqlite.create_pool("tokens.db") as db_pool:
            async with IvanTheRagingBot(
                package_location,
                db_pool=db_pool,
                settings=settings,
            ) as bot:
                await bot.start()

    try:
        asyncio.run(runner())
    except KeyboardInterrupt:
        LOGGER.warning("Shutting down due to Keyboard Interrupt...")


if __name__ == "__main__":
    main()
