from dotenv import load_dotenv

from ivantheraginbot.bot import ChatReader


def main():
    load_dotenv()
    bot = ChatReader()
    bot.run()


if __name__ == "__main__":
    main()
