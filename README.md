# PDFBanaLeBro

[![Add Me](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/api/oauth2/authorize?client_id=933041693970800691&permissions=8&scope=bot)

PDFBanaLeBro is a discord bot using for PDF utility.
This bot was made because some professors were too lazy to send PDF notes

The purpose of this bot is to convert sent images to PDF and send them to discord chat.

* Commands start with PreCommand (set to "soja." in modules)

The Bot currently has following commands.
1. soja.ping - Sends ping result
2. soja.start - Initiates a session
3. soja.f - Should be added as caption to every image you want as PDF
4. soja.end - Ends the active session and sends PDF to target chat (Doesn't work no images has been sent or without any active session)
5. soja.debug - placeholder command for testing purpose (disabled when app_mode is set to false in config.ini)

Although usable, the bot lacks some features as of now. Below are the planned features (Todo list)
1. Support for encrypted PDFs
2. Extract images from a PDF and send them as a zip
3. A new soja.help command to explain command usage
4. Support for multiple images in one message (The bot was designed for circumstances when you will only send one image at a time so not implemented by default)
5. Improve image fetching mechanism (At present, bot awaits soja.end command to convert every images to PDF. A better mechanism is to append images to existing PDF)
6. soja.end doesn't work when images sent are zero. So introduce a command to end session forcefully

So, How should you use the bot?
- Simple add the bot to your server using the button provided above
- Start a session using ```soja.start```
- Send all your images one by one each with caption ```soja.f```
- Send ```soja.end``` to end the session and receive the final PDF
- And your'e done

Pull requests are welcomed (unless you're a spam bot)