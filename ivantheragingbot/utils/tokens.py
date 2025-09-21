import asyncio

import twitchio

from ivantheragingbot.config import Settings


async def main() -> None:
    settings = Settings()

    async with twitchio.Client(
        client_id=settings.CLIENT_ID, client_secret=settings.CLIENT_SECRET
    ) as client:
        await client.login()
        user = await client.fetch_users(
            logins=["ivantheragingpython", "ivantheraginbot"]
        )
        for u in user:
            print(f"User: {u.name} - ID: {u.id}")


if __name__ == "__main__":
    asyncio.run(main())
