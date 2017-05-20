---
layout: post
title: "ChatOps on Windows with Hubot and PowerShell"
date: 2016-01-25 13:37:00
comments: true
description: Using Slack, Hubot and PowerShell to enable ChatOps in the Microsoft Ecosystem.
modified: 2016-07-08
---

* **UPDATED** 20th May, 2017: If you are more comfortable with Python, you can [this post]({% post_url 2017-05-29-chatops-with-powershell-and-errbot %}) on Errbot and PowerShell.
* **UPDATED** 8th July, 2016: Created a installation video for Hubot using PowerShell DSC [here](https://www.youtube.com/watch?v=Gh-vYprIo7c).
* **UPDATED** 1st July, 2016: Created a PowerShell DSC Resource to install Hubot which makes the process much easier. Can be found on GitHub [here](https://github.com/MattHodge/Hubot-DSC-Resource) or found in the PowerShell Gallery [here](https://www.powershellgallery.com/packages/Hubot).

ChatOps is a term used to describe bringing development or operations work that is already happening in the background into a common chat room. It involves having everyone in the team in a single chat room, then bringing tools into the room so everyone can automate, collaborate and see how automation is used to solve problems. In doing so, you are unifying the communication about what work gets done and have a history of it happening.

ChatOps can be supplemented with the use of tools or scripts exposed using a chat bot. Users in the chat room can talk to the bot and have it take actions on their behalf, some examples of this may be:

* Checking the status of a Windows Service
* Finding out who is on call via the PagerDuty API
* Querying a server to see how much disk space is available

Bots can also be a great way to expose functionality to low-privledged users such as help desk staff, without having to create web interfaces or forms.

If you want more details on the concept of ChatOps, I recommend watching **[ChatOps, a Beginners Guide](https://www.youtube.com/watch?v=F8Vfoz7GeHw)** presented by [Jason Hand](https://twitter.com/jasonhand).

A popular toolset for ChatOps is [Slack](https://slack.com/) as the chat client, and [Hubot](https://hubot.github.com/) as the bot. In this post we will use Slack and Hubot together with a PowerShell module I've written called [PoshHubot](https://github.com/MattHodge/PoshHubot). The module will handle installation and basic administration of Hubot. From there, we will integrate Hubot with PowerShell so we can perform some ChatOps in the Microsoft ecosystem.

* TOC
{:toc}

## Basic Hubot Concepts

There are a few basic Hubot concepts I want to introduce to you before we continue.

### Node.js and CoffeeScript
Hubot is built in CoffeeScript, which is a programming language that complies into JavaScript. Hubot is built on Node.js. This means the server running the bot will need to have Node.js and CoffeeScript installed. The [PoshHubot](https://github.com/MattHodge/PoshHubot) module will handle this.

When writing scripts for your bot, you will have to get your hands a little dirty with CoffeeScript. We will be calling PowerShell from inside CoffeeScript, so we only need to know a tiny bit to get by.

### Environment Variables
Hubot scripts make heavy use of environment variables to set certain options for the bot.

One example of this is to allow the Hubot to access sites with invalid SSL certificates, you would set an environment variable of `NODE_TLS_REJECT_UNAUTHORIZED`.

There are 3 possible ways to do this:

* You can set these environment variables as a system wide setting in an Administrative PowerShell prompt using:

{% gist 9fe3f9bee81c2d6f5327 %}

* You can set them in the current PowerShell instance before you start the bot using:
{% highlight powershell %}
$env:NODE_TLS_REJECT_UNAUTHORIZED = '0'
{% endhighlight %}

* You can store the environment variable in the `config.json` file that we generate during the Hubot installation, which the `PoshHubot` module will load before starting the bot.


### Bot Brain
Hubot has a *brain*, which is simply a place to store persist data. For example, you could write a script to have Hubot store URL's for certain services, which you could append to via chat commands. You want these URL's to persist after Hubot reboots, so it needs to save them to its brain.

There are many brain adapters for Hubot, for example MySQL, Redis and Azure Blob Storage. For this blog we install a file brain - which will just store the brain as a *.json* file on the disk.

## Requirements
You will need to a have a few things ready to get a Hubot setup with Slack:

1. A Windows Machine with PowerShell 4.0+. For this tutorial I will be using a Windows 2012 R2 Standard machine with GUI. Once you get comfortable with Hubot you may decide to switch to using Server Core, which is a great choose for running Hubot
1. Administrative access in your Slack group to create a Hubot integration

## Create a Slack Integration for Hubot

To have Hubot communicating with Slack, we need to configure an integration. From Slack:

* Choose **Apps & Custom Integrations**

![Slack - Apps & Custom Integrations](/images/posts/chatops_on_windows/slack_apps_customize.png "Slack - Apps & Custom Integrations")

* Search for **Hubot** and choose **Install**
* Provide a Hubot username - this will be the name of your bot. For this blog the bot will be called bender
* Click **Add Hubot Integration**

After the integration has been added, you will be provided an API Token, something like `xoxb-XXXXX-XXXXXX`. We will need this later so note it down.

Additionally, you can customize your bots icon and add channels at this screen.

![Slack - Bot name & Icon](/images/posts/chatops_on_windows/slack_choose_icon.png "Slack - Bot name & Icon")

## Installing Hubot
Install the `PoshHubot` module in one of two ways:

* Install it using PowerShell get with `Install-Module -Name PoshHubot`
* Manually download the [latest PoshHubot release](https://github.com/MattHodge/PoshHubot/releases) and extract the module your PowerShell Module directory.

Once installed, we are going to create a configuration file that `PoshHubot` will use.

{% highlight powershell %}
# Import the module
Import-Module -Name PoshHubot -Force

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
* [Forever](https://github.com/foreverjs/forever) which will run Hubot as a background process

### Removing Hubot Scripts

Hubot comes installed with some default scripts which are not required when running on Windows. We can use the `Remove-HubotScript` command to remove them.

{% highlight powershell %}
# Be sure to provide the correct ConfigPath for your bot
Remove-HubotScript -Name 'hubot-redis-brain' -ConfigPath 'C:\PoshHubot\config.json'
Remove-HubotScript -Name 'hubot-heroku-keepalive' -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}


### Installing Hubot Scripts

There are many third party scripts available for Hubot to extend its functionality. You can search for them on the [Node.js package manager site](https://www.npmjs.com/search?q=hubot) or on [GitHub](https://github.com/hubot-scripts). We will use the `Install-HubotScript` function to install some useful scripts.

{% highlight powershell %}
# Authentication Script, allowing you to give permissions for users to run certain scripts
Install-HubotScript -Name 'hubot-auth' -ConfigPath 'C:\PoshHubot\config.json'
# Allows reloading Hubot scripts without having to restart Hubot
Install-HubotScript -Name 'hubot-reload-scripts' -ConfigPath 'C:\PoshHubot\config.json'
# Stores the Hubot brain as a file on disk
Install-HuBotScript -Name 'jobot-brain-file' -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

## Starting Hubot

Before we can start our bot and connect it to Slack, we have to configure the environment variables required by the scripts we are using. A good way to find out what environment variables a script is using is to look it up on GitHub. For example,
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

With all the configuration in place, we can start Hubot.

{% highlight powershell %}
Start-Hubot -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

Open up Slack and see your bot online! If for some reason your bot doesn't connect, you can find the logs in the `LogPath` defined earlier in the `config.json` file.

Hubot comes with some built in commands, so you can directly message your bot with `help` and see if you get a response back.

![Speaking to Hubot for the first time](/images/posts/chatops_on_windows/speaking_to_hubot_in_slack.png "Speaking to Hubot for the first time")

If you want your bot to join certain channels, you can enter `/invite @bender` in Slack to bring him into the channel. To have Hubot perform commands, you need to address him in the channel. Try a `@bender pug bomb me`.

![Yay! Pug Bombed!](/images/posts/chatops_on_windows/bot_pug_bomb.png "Yay! Pug Bombed!")

## Integrating Hubot with PowerShell

We have our Hubot joined to Slack and we have triggered a few pug bombs, but it is time to do something useful - create our own script.

The [Hubot documentation](https://hubot.github.com/docs/scripting/) covers scripting in detail and I recommend giving it a read before continuing on.

We are going to write a basic script to find the status of a Windows service on the machine hosting the Hubot. The plan is:

* Send the bot a message saying `@bender: get service dhcp` - where `dhcp` could be the name of any service.
* The Hubot script will use a regex capture group to select out the name of the service (in this case `DCHP`)
* The Hubot script will pass this captured service name into a PowerShell script to find the status of the service.
  * If the service exists, it will return the service status.
  * If the service does not exist, it will say so.
* The PowerShell script will return the results in a json format. This will make it far easier to work with in CoffeeScript

### Install Edge.js and Edge-PS

[Edge.js](https://github.com/tjanczuk/edge) and [Edge-PS](https://github.com/dfinke/edge-ps) are Node.js packages which allow calling .NET and PowerShell (among other things) from Node.js.

To use them inside Hubot, we need to add them to the `package.json` file which is generated when we install Hubot for the first time. You can find `package.json` in the `BotPath` specified above, in our case it is `C:\myhubot\packages.json`. We will also add a version constraint. The latest version of each package can be found by searching the [npm package manager](https://www.npmjs.com).

After you have added them your `package.json` should look similar to this:

{% highlight json %}
{
  "name": "bender",
  "version": "0.0.0",
  "private": true,
  "author": "PoshHubot <posh@hubot.com>",
  "description": "PoshHubot is awesome.",
  "dependencies": {
    "hubot": "^2.18.0",
    "hubot-diagnostics": "0.0.1",
    "hubot-google-images": "^0.2.6",
    "hubot-google-translate": "^0.2.0",
    "hubot-help": "^0.1.3",
    "hubot-heroku-keepalive": "^1.0.2",
    "hubot-maps": "0.0.2",
    "hubot-pugme": "^0.1.0",
    "hubot-redis-brain": "0.0.3",
    "hubot-rules": "^0.1.1",
    "hubot-scripts": "^2.16.2",
    "hubot-shipit": "^0.2.0",
    "hubot-slack": "^3.4.2",
    "edge": "^5.0.0",
    "edge-ps": "^0.1.0-pre"
  },
  "engines": {
    "node": "0.10.x"
  }
}
{% endhighlight %}

Usually, after you update the `package.json`, you would need to run  `npm install` to download the packages that have been added. This is handled behind the scenes for you by the `Start-Hubot` command.

### Create the PowerShell script

We need to design a PowerShell script that can be called from CoffeeScript, the Hubot scripting language. I recommend using the following methods when creating PowerShell scripts that will be called from Hubot:.

* **Create PowerShell functions with paramaters** - It makes it nice and easy to call PowerShell from CoffeeScript when they have well defined parameters
* **Use error handling in your script** - Use try-catch blocks inside your PowerShell functions so you can return a message to the bot if the command has failed
* **Output results in json** - This is a great way to pass data back to CoffeeScript. You can use PowerShell objects to send data back and have CoffeeScript pick out the parts you want

Keeping these methods in mind, I created a `Get-ServiceHubot` function to find the service status:

{% gist 66bf00bedb98d72c2506 %}

I am applying some [Slack formatting](https://get.slack.help/hc/en-us/articles/202288908-Formatting-your-messages) in my output, including the use of asterisks around words for bold and back ticks for code blocks. You will notice there are double backticks in the code so PowerShell does not interpret them.

Here is some example output from the PowerShell when the function is run against a service that exists:

{% highlight powershell %}
# Dot Source the function
. .\Get-ServiceHubot.ps1

# Get a service that exists on the system
Get-ServiceHubot -Name dhcp
{% endhighlight %}

{% highlight json %}
{
  "success":  true,
  "output":  "Service dhcp (*DHCP Client*) is currently `Running`"
}
{% endhighlight %}

Here is some example output from the PowerShell when the function is run against a service that doesn't exist on the machine:

{% highlight powershell %}
# Get a service that exists on the system
Get-ServiceHubot -Name MyFakeService
{% endhighlight %}

{% highlight json %}
{
  "success":  false,
  "output":  "Service MyFakeService does not exist on this server."
}
{% endhighlight %}

Save the PowerShell function into the `scripts` folder in the Hubot directory. In my case I will be saving it to `C:\myhubot\scripts\Get-ServiceHubot.ps1`

### Create the Hubot script

Now that our PowerShell function is completed, we need to wire it up to Hubot using CoffeeScript.

The goal for the CoffeeScript portion is to take a users message to the bot, work out the service name, pass it into the PowerShell script and return the result to the user.

This is the script I designed to call the PowerShell function. Be sure to read the comments so you understand how it works.

{% gist 8e3461622b7f3c1f9a53 %}

Save the CoffeeScript into the Hubot scripts directory as well, in my case this will be `C:\myhubot\scripts\get-servicehubot.coffee`.

### Testing the script

To load the script into Hubot, you need to restart the bot:

{% highlight powershell %}
Restart-Hubot -ConfigPath 'C:\PoshHubot\config.json'
{% endhighlight %}

You will notice that `npm` installs the Edge.js dependencies we added in the `package.json`.

![Automatic installation of npm dependencies](/images/posts/chatops_on_windows/start_hubot_install_npm_packages.png "Automatic installation of npm dependencies")

Check the logs in `LogPath` defined earlier in the `config.json` file to make sure that Hubot started successfully and loaded your script. You should see a line in the log file like this:
{% highlight text %}
[Sun Jan 24 2016 10:25:59 GMT-0800 (Pacific Standard Time)] DEBUG Parsing help for C:\myhubot\scripts\get-servicehubot.coffee
{% endhighlight %}

When your bot joins the channel, ask it for help again. You will notice that the `get service` command has been added to the help. This is done automatically when you fill out the header part of the CoffeeScript script.

![Command added to Hubot Help](/images/posts/chatops_on_windows/hubot-help-from-comments.png "Command added to Hubot Help")

Now you can try some `get service <service>` commands and see the results:

![Get service is now completed!](/images/posts/chatops_on_windows/get-service-results.png "Get service is now completed!")

## Wrapping Up

Hubot with PowerShell is a fantastic way to bring automation to your environment. With a tiny amount of CoffeeScript you can take your pre-existing PowerShell functions and make them available in a chat channel for anyone in your team to access. This is especially useful for allowing people in your company to access information on-demand from places Operations teams may only have access to.

I'd love to hear about the cool scripts you come up with when leveraging Hubot and PowerShell! Tweet me [@matthodge](https://twitter.com/matthodge).
