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

A popular combination of tools for ChatOps is [Slack](https://slack.com/) for the chat client, and [Hubot](https://hubot.github.com/) as the bot, which is what this post will be targeting.

* TOC
{:toc}

## Requirements
You will need to a have a few things ready to get a Hubot setup with Slack:

1. A Windows Machine with PowerShell 4.0+. For this tutorial I will be using a Windows 2012 R2 Standard machine with GUI. Once you get comfortable with Hubot you may decide to switch to using Server Core which is a great choose for running Hubot
1. Slack Administrative access to create a Hubot integration, which we will cover below

## Create a Slack Integration for Hubot

To have Hubot speaking to Slack, we need to configure an integration. From Slack:
* Choose **Apps & Custom Integrations**

![Slack - Apps & Custom Integrations](/images/posts/chatops_on_windows/slack_apps_customize.png "Slack - Apps & Custom Integrations")

* Search for **Hubot** and choose **Install**
* Provide a Hubot username - this will be the name of your bot. For this blog the bot will be called *bender*. Click **Add Hubot Integration**

After the integration has been added, you will be provided an API Token, something like `xoxb-XXXXX-XXXXXX`. We will need this later so note it down.

Additionally, you can customize your bots icon and add channels at this screen.

![Slack - Bot name & Icon](/images/posts/chatops_on_windows/slack_choose_icon.png "Slack - Bot name & Icon")

## Basic Hubot Concepts

There are a few basic Hubot concepts I want to introduce to you before we continue.

### Node.js and CoffeeScript
Hubot is built in CoffeeScript, which is a programming language that complies into JavaScript. Hubot is built on Node.js. This means the server running the bot will need to have Node.js and CoffeeScript installed, but don't worry, my PowerShell function will handle all this for you.

When writing scripts for your bot, you will have to get your hands a little dirty with CoffeeScript. We will have it calling PowerShell scripts, so we only need to know a tiny bit of CoffeeScript to get by.

### Environment Variables
Hubot and its addons / scripts makes heavy use of environment variables to set certain options for the bot.

One example of this is for the Slack adapter, it expects the `HUBOT_SLACK_TOKEN` environment variable to be set. The easiest way to set these is with PowerShell

You can set these environment variables as a system wide setting using:
{% highlight powershell %}
# This will need to be done with an Administrative PowerShell Prompt
[Environment]::SetEnvironmentVariable("HUBOT_SLACK_TOKEN", "xoxb-XXXXX-XXXXXX", "Machine")
{% endhighlight %}

Or you can set them in the current PowerShell instance before you start the bot  using:
{% highlight powershell %}
  $env:HUBOT_SLACK_TOKEN = 'xoxb-XXXXX-XXXXXX'
{% endhighlight %}

## Installing Hubot
Jekyll also offers powerful support for code snippets:

{% highlight ruby %}
def print_hi(name)
  puts "Hi, #{name}"
end
print_hi('Tom')
#=> prints 'Hi, Tom' to STDOUT.
{% endhighlight %}

## Connecting Hubot to Slack

## Integrating Hubot with PowerShell
