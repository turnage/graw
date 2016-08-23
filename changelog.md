#### 0.4.0 August 23rd, 2016

Added two public packages:

graw/grerr defines values for common errors so bots can define scenarios for
them.

graw/botfaces defines the interfaces graw looks for in bots so users can type
assert their bots and verify their own implementations.

#### 0.3.0 August 21st, 2016

Removed ````Scrape()```` from api.

#### 0.2.2 August 11th, 2016

Fixed [issue](https://github.com/turnage/graw/issues/13) reported by [silviucm](https://github.com/silviucm).

#### 0.2.1 June 29th, 2016

Fixed [issue](https://github.com/turnage/graw/issues/12) reported by [jonas747](https://github.com/jonas747).

graw now uses the bot's custom user agent when initiating OAuth2 relationships with Reddit so graw bots are not
subject to rate limiting they didn't earn.

#### 0.2.0 June 12th, 2016

Accepted [pr](https://github.com/turnage/graw/pull/8) from [ultralight-meme](https://github.com/ultralight-meme).

Allows bots to set a custom refresh time.

#### 0.1.0 May 25th, 2016

Accepted [pr](https://github.com/turnage/graw/pull/9) from [ultralight-meme](https://github.com/ultralight-meme).

Requests raw json from reddit so stylized elements in text such as bold come in source text format.

For example, a message: "containing a **bold** word" yields ```containing a **bold** word``` instead of
``containing a <b>word</b> word```.
