from dotenv import load_dotenv

from ivantheragingbot.bot import ChatReader

if __name__ == "__main__":
    load_dotenv()
    bot = ChatReader()
    bot.run()
