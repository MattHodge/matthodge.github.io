---
layout: post
title:  "ChatsOps on Windows with Hubot and PowerShell"
date:   2016-01-17 13:37:00
comments: false
---

ChatOps is a term used to describe conversation driver development or operations for an Ops or Development team. It involves having everyone in the teams in a single chatroom, and brining tools into the chatroom which can help the team to automate, collaborate and work better as a team. The tools or automations are usually exposed using a chat bot which users in the chatroom can talk to to have it take actions, some examples of this may be:

* Have the bot kick off a script
* Check who is on call via the PagerDuty API
* Query a server to see how much disk space is available

Bots can also be a great way to expose functionality to low-privledged users such as help desk staff, without having to create web UI's or forms.

I won't go into any details on the concepts of ChatOps, but I recommend watching **[ChatOps, a Beginners Guide](https://www.youtube.com/watch?v=F8Vfoz7GeHw)** presented by [Jason Hand](https://twitter.com/jasonhand) if you are new to the term.

A popular combination of tools for ChatOps is [Slack](https://slack.com/) for the chat client, and [Hubot](https://hubot.github.com/) as the bot, which is what this post will be targeting. This post will also be using a PowerShell module I've written called [Hubot-PowerShell](https://github.com/MattHodge/Hubot-PowerShell). The module will handle installation and basic administration Hubot.

* TOC
{:toc}

## Basic Hubot Concepts

There are a few basic Hubot concepts I want to introduce to you before we continue.

### Node.js and CoffeeScript
Hubot is built in CoffeeScript, which is a programming language that complies into JavaScript. Hubot is built on Node.js. This means the server running the bot will need to have Node.js and CoffeeScript installed. The `Hubot-PowerShell` module will handle this.

When writing scripts for your bot, you will have to get your hands a little dirty with CoffeeScript. We will be calling PowerShell from inside CoffeeScript, so we only need to know a tiny bit to get by.

### Environment Variables
Hubot and its addons / scripts makes heavy use of environment variables to set certain options for the bot.

One example of this is to allow the Hubot to access sites with invalid SSL certificates, you would set an environment variable of `NODE_TLS_REJECT_UNAUTHORIZED`.

There are 3 possible ways to do this:

* You can set these environment variables as a system wide setting in an Administrative PowerShell prompt using:
{% highlight powershell %}
# This will need to be done with an Administrative PowerShell Prompt
[Environment]::SetEnvironmentVariable("NODE_TLS_REJECT_UNAUTHORIZED", "0", "Machine")
{% endhighlight %}

* You can set them in the current PowerShell instance before you start the bot  using:
{% highlight powershell %}
  $env:NODE_TLS_REJECT_UNAUTHORIZED = '0'
{% endhighlight %}

* You can store the environment variable in the `config.json` file that we generate during the Hubot installation, which the  `Hubot-PowerShell` module will load before starting the bot.


### Bot Brain
Hubot has a *brain* which is simply a place to store data you want to persist after Hubot reboots. For example, you could write a script to have Hubot store URL's for certain services, which you could append to via chat commands. You want these URL's to persist after Hubot reboots, so it needs to save them to its brain.

There are many brain adapters for Hubot, for example MySQL, Redis and Azure Blob Storage. For this blog I will be using a file brain - which will just store the brain as a *.json* file on the disk.

## Requirements
You will need to a have a few things ready to get a Hubot setup with Slack:

1. A Windows Machine with PowerShell 4.0+. For this tutorial I will be using a Windows 2012 R2 Standard machine with GUI. Once you get comfortable with Hubot you may decide to switch to using Server Core which is a great choose for running Hubot
1. Administrative access in your Slack group to create a Hubot integration

## Create a Slack Integration for Hubot

To have Hubot speaking to Slack, we need to configure an integration. From Slack:

* Choose **Apps & Custom Integrations**

![Slack - Apps & Custom Integrations](/images/posts/chatops_on_windows/slack_apps_customize.png "Slack - Apps & Custom Integrations")

* Search for **Hubot** and choose **Install**
* Provide a Hubot username - this will be the name of your bot. For this blog the bot will be called *bender*
* Click **Add Hubot Integration**

After the integration has been added, you will be provided an API Token, something like `xoxb-XXXXX-XXXXXX`. We will need this later so note it down.

Additionally, you can customize your bots icon and add channels at this screen.

![Slack - Bot name & Icon](/images/posts/chatops_on_windows/slack_choose_icon.png "Slack - Bot name & Icon")

## Installing Hubot
Install the `Hubot-PowerShell` module by downloading it from git and placing it into your PowerShell Module directory.

First we are going to create a configuration file that `Hubot-PowerShell` will use.

{% highlight powershell %}
# Import the module
Import-Module -Name Hubot-PowerShell -Force

# Create hash of configuration options
$newBot = @{
    Path = "C:\PoshHubot\config.json"
    BotName = 'bender'
    BotPath = 'C:\myhubot'
    BotAdapter = 'slack'
    BotOwner = 'Matt <matt@email.com>'
    BotDescription = 'my awesome bot'
    LogPath = 'C:\PoshHubot\Logs'
    BotDebugLog = $true
}

# Splat the hash to the CmdLet
New-PoshHubotConfiguration @newBot
{% endhighlight %}

Next, we need to install all the required components for Hubot, which will be handled by the `Install-Hubot` command.

{% highlight powershell %}
# Install Hubot
Install-Hubot -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

This will install the following:

* Chocolatey
* Node.js
* Git
* CoffeeScript
* Hubot Generator
* [Forever](https://github.com/foreverjs/forever) which will run Hubot as a background process.

### Removing Hubot Scripts

Hubot comes installed with some default scripts which are not required when running on Windows. We can use the `Remove-HubotScript` command to remove them.

{% highlight powershell %}
# Be sure to provide the correct ConfigPath for your bot
Remove-HubotScript -Name 'hubot-redis-brain' -ConfigPath 'C:\PoshHubot\config.json'
Remove-HubotScript -Name 'hubot-heroku-keepalive' -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}


### Installing Hubot Scripts

We will now install some third party Hubot scripts using `Install-HubotScript`.

{% highlight powershell %}
# Authentication Script, allowing you to give permissions for users to run certain scripts
Install-HubotScript -Name 'hubot-auth' -ConfigPath 'C:\PoshHubot\config.json'
# Allows reloading Hubot scripts without having to restart Hubot
Install-HubotScript -Name 'hubot-reload-scripts' -ConfigPath 'C:\PoshHubot\config.json'
# Stores the Hubot brain as a file on disk
Install-HuBotScript -Name 'jobot-brain-file' -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

## Starting Hubot

Before we can start our bot and connect it to slack, we have to configure the environment variables required by the scripts we are using. A good way to find out what environment variables a script is using is to look it up on github. For example,
the [jubot-brain-file](https://github.com/8DTechnologies/jobot-brain-file) script requires `FILE_BRAIN_PATH` to be set.

![jubot-brain-file Script](/images/posts/chatops_on_windows/bot_brain_coffee.png "jubot-brain-file Script")

Additionally, the [Slack adapter](https://github.com/slackhq/hubot-slack) for Hubot requires an environment variable to be set for the Slack API token called `HUBOT_SLACK_TOKEN`.

We will store both of these in the `config.json` file we created earlier.

Open the `C:\PoshHubot\config.json` file and in the `EnvironmentVariables` section, add the new environment variables.

The completed `config.json` file should look something like this:

{% highlight json %}
{
  "Path": "C:\\PoshHubot\\config.json",
  "BotAdapter": "slack",
  "BotDebugLog": {
    "IsPresent": true
  },
  "BotDescription": "my awesome bot",
  "BotPath": "C:\\myhubot",
  "BotOwner": "Matt <matt@email.com>",
  "LogPath": "C:\\PoshHubot\\Logs",
  "BotName": "bender",
  "ArgumentList": "--adapter slack",
  "BotExternalScriptsPath": "C:\\myhubot\\external-scripts.json",
  "PidPath": "C:\\myhubot\\bender.pid",
  "EnvironmentVariables": {
    "HUBOT_ADAPTER": "slack",
    "HUBOT_LOG_LEVEL": "debug",
    "HUBOT_SLACK_TOKEN": "xoxb-XXXXX-XXXXXX",
    "FILE_BRAIN_PATH": "C:\\PoshHubot\\"
  }
}
{% endhighlight %}

With all the configuration in place, we can start Hubot!

{% highlight powershell %}
Start-Hubot -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

Open up Slack and check if your bot came online! Hubot comes with some built in commands, so you can directly message your bot with `help` and see if you get a response back.

![Speaking to Hubot for the first time](/images/posts/chatops_on_windows/speaking_to_hubot_in_slack.png "Speaking to Hubot for the first time")


## Integrating Hubot with PowerShell
