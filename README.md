# ivantheragingbot

Golang rewrite of the Twitch TTS bot. It joins a channel, listens to chat, and speaks messages or plays bundled alerts.

## Features
- Google translate TTS (Spanish by default) for commands and optional automatic chat reading.
- Commands: `!speak`, `!help`, `!red`, `!alerta`, `!categoria`.
- Ignores the bot user, Nightbot, and StreamElements unless `READ_AUTHOR_MESSAGE` is set.
- Plays local MP3 assets from `assets/` for the alert-style commands.

## Setup
1. Install Go 1.25+ and ensure audio playback works on the host.
2. Copy `example.env` to `.env` and fill in:
   - `TWITCH_USERNAME` (bot account username)
   - `IRC_TOKEN` (Twitch IRC OAuth token, usually starts with `oauth:`)
   - `TWITCH_CHANNEL` (channel to join)
   - Optional flags: set `AUTO_READ` to have regular chat messages spoken; set `READ_AUTHOR_MESSAGE` to allow ignored users to be read; override language with `TTS_LANG`/`TTS_TLD`.
3. Run the bot: `go run .`

The bot uses the default system audio output. TTS generation and Twitch connectivity require internet access when running.
