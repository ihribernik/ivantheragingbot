import logging


class ErrorHandler:
    # def __init__(self):
    #     self.logger = logging.getLogger("ErrorHandler")

    async def handle_error(self, ctx, error):
        """Captura y registra errores."""
        self.logger.error("Error en el comando %s:%s", ctx.command, error)
        await ctx.send("Ocurrió un error. Inténtalo de nuevo más tarde.")
