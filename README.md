# rss-tg-chan
A service that converts a bunch of RSS feeds into a Telegram channel.


## Overview

* Uses JSON as a source of RSS feeds (you can put it on [GitHub Gist](https://gist.github.com))
* Uses JSON on GitHub Gist as a storage for the state (just for fun C:)
* Publishes posts from those RSS feeds to a Telegram channel


## Quick start

### 1. Create your feeds source
1. Create a JSON file with the following structure and put it somewhere: 
```json
{
  "feeds": [
    "https://blog.github.com/all.atom",
    "https://github.blog/engineering.atom",
    "http://www.reddit.com/r/golang/.rss"
  ]
}
```
2. If you want your bot to see new RSS feeds without restart, make sure your link always points to the latest version of the file (e.g. `https://gist.githubusercontent.com/rozag/dcd1b09bbe12a942dbe0f3bcbb2ace7b/raw/feeds-tech.json`). See instructions for GitHub Gist [here](https://stackoverflow.com/a/47175630)

### 2. Create GitHub Gist storage JSON and token
1. [Create](https://gist.github.com/) a new private GitHub Gist with the following content: `{}`. Remember the `.json` file name and the Gist id (you can find it in the URL)
2. [Create](https://github.com/settings/tokens/new) a GitHub token with `gist` scope. Remember the generated token

### 3. Create Telegram bot and channel
1. Create a Telegram channel. Remember it's id (e.g. @my_channel)
2. Ask [BotFather](https://t.me/botfather) to create a new bot for you. Remember your bot token
3. Add your new bot to administrators of your new channel

### 4. Build and deploy the container
1. BTW you can simply run the app with `go run .`
2. Clone this repo `git clone https://github.com/rozag/rss-tg-chan.git`
3. `cd rss-tg-chan`
4. Put all your ids and tokens to the `config.ini` file like this:
```
source=LINK_TO_YOUR_SOURCES_JSON

githubToken=YOUR_GITHUB_TOKEN_WITH_GIST_SCOPE
githubGistID=YOUR_STORAGE_JSON_GIST_ID
githubGistFileName=YOUR_STORAGE_JSON_GIST_FILE_NAME

tgToken=YOUR_TELEGRAM_BOT_TOKEN
tgChannel=@YOUR_TELEGRAM_CHANNEL_ID
```
5. Build docker image: `docker build --rm --tag YOUR_IMAGE_NAME:YOUR_IMAGE_TAG .`
6. Run your container locally `docker run YOUR_IMAGE_NAME:YOUR_IMAGE_TAG` or use [`docker save`](https://docs.docker.com/engine/reference/commandline/save/), [`scp`](https://unix.stackexchange.com/a/106482) and [`docker load`](https://docs.docker.com/engine/reference/commandline/load/) to transfer it to some remote host


## Contributing

If you find a bug - [create an issue](https://github.com/rozag/rss-tg-chan/issues/new). It's your contribution. And PRs are always welcome.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
