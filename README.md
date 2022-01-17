# telegrambot
General Pub/Sub Telegram Bot

What is the idea here: a simple bot for telegram. It allows you to send notifications via http and their a brodcast to whoever subscribes a channel.

Users should login 1st, so the bot will not send data to strangers

How it works: 
- Bot will persist sessions and subscriptions on postgres database.
- An APIKey is set when setting up the service
- Http request should be sent w X-API-KEY header set, see below.
- Users call bot
- User call /login <pass> where <pass> is the password
- User can call:
    - /sub <chan> to subscribe a channel
    - /unsub <chan> to unsubscribe
    - /mysubs to list user own subscriptions
    - /subbers <chan> to list all users subscribing a channel
    - /pub <chan> <msg> will send a message to a channel


In case http requests are received as shown below, they will be forwarded to telegram users

## Env

These vars are required to run the service:

- TBOTKEY: Telegram API Key for the bot
- DSN postgres database dsn, compatible with pg and gorm
- APIKEY: This is the apikey http clients should use when sending a notification
- USERPASS: Password users should give to bot when starting to comunicate

Please refer to:
[Telegram Bot Notify](https://github.com/digitalcircle-com-br/telegrambot-notify)

for a simple notification lib to be used as companion to this bot.

## Sample call
```http
POST https://tbot.digitalcircle.com.br/pub
X-API-KEY: {{xapikey}}

{
    "ch":"a",
    "msg":"FAFA 💀 from api"
}
```