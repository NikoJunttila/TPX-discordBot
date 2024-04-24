# Discord-Twitch Bot

This Discord bot allows you to merge Twitch chat with a Discord channel.<br>
Checks RIOT API for new matches for certain users and shows stats when game ends and check users latest matches for their stats, check live game and more functionalities.<br>
It's built using Golang and utilizes the Riot Games API for certain functionalities and the Twitch API for integrating Twitch chat. Additionally, it stores relevant data in a PostgreSQL database.

## Features

- **Twitch Chat Integration**: The bot seamlessly merges Twitch chat messages into a specified Discord channel, providing a unified experience for your community.
  
- **Riot Games API Integration**: Utilizes the Riot Games API to access data or functionalities related to Riot Games, enhancing the bot's capabilities for gaming communities.

- **PostgreSQL Database**: Stores necessary data in a PostgreSQL database for efficient data management and retrieval.

## Installation

1. **Clone the Repository:**
```
git clone https://github.com/NikoJunttila/TPX-discordBot
```
2. **Install Dependencies:**
```
go mod tidy
```
3. **Set Up Environment Variables:**
create .env file. copy example.env variables there and fill them with your own apikeys/discord channels ids
4. **Run the Bot:**
go build -o discordbot //change discordbot to anything you want to.<br>
execute discordbot
## Usage

- Invite the bot to your Discord server and ensure it has appropriate permissions to read and send messages in the desired channels.
