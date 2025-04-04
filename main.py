from pathlib import Path
from dotenv import load_dotenv

from ivantheraginbot.bot import ChatReader


def main():
    load_dotenv()
    package_location = Path.cwd()
    bot = ChatReader(package_location)
    bot.run()


if __name__ == "__main__":
    main()
