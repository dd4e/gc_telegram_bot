# gc_telegram_bot

## About
Garbage collector telegram bot.
Designed to remove outdated messages from all members of the group. 
The bot only saves the message metadata. Radis is used as a backend for storing messages.

## How-To

### Installation
```
go get github.com/dadmoscow/gc_telegram_bot
go build gc_telegram_bot
```

### Configuration
#### System environment
**GC_TOKEN**  
Telegram bot TOKEN  

**GC_CHECK_TIMEOUT**  
Timeout checking for old messages in seconds  
*Default:* 60 sec

**GC_TIMEOUT_LIMIT**  
The maximum time limit for storing messages in seconds   
*Default:* 604800 sec

**GC_REDIS_ADDR**  
Redis address in format *ip*:*port*  
*Default:* "127.0.0.1:6379"

**GC_REDIS_DB**  
Radis database number  
*Default:* 0

**GC_REDIS_PWD**  
Database password  
*Default*: None

**GC_USE_SOCKS5**  
Use SOCKS5 proxy to connect  
*Default*: false

**GC_SOCKS5_ADDR**  
SOCKS5 address in format *ip*:*port*  

**GC_SOCKS5_USER**  
SOCKS5 username  

**GC_SOCKS5_PWD**  
SOCKS5 password  

**GC_BOT_DEBUG**  
Debug mode  
*Default*: false  

#### Configuration file
To-Do

### Running

## Bot Commands
/help 		-- print this message  
/on   		-- the bot will delete outdated messages  
/off		-- the bot will be disabled  
/timeout	-- new timeout after which the messages will be deleted  
/delete		-- delete all messages  
/setting	-- print current settings  
/stop		-- !!! Delete all messages, delete settings and stop the bot !!!  

## To-Do List
