package app

const HelpMsg = `
Supported command:
/help 		-- print this message
/on   		-- the bot will delete outdated messages
/off		-- the bot will be disabled
/timeout	-- new timeout after which the messages will be deleted
/delete		-- delete all messages
/setting	-- print current settings
/stop		-- !!! DeleteMessage all messages, delete settings and stop the bot !!!

Timeout format:
Timeout is set in the format: <decimal><unit suffix>
unit suffix one of "s", "m", "h"
Example: 1h15m, 24h, 30m, 60s, 10h30m15s
`

const StartMsg = `
Instructions to get started:
1. Add bot to group
2. Give him admin rights to delete messages
3. In the group send a command to the bot /on to run
4. You can change the timeout setting
5. To get the list command send /help
`
