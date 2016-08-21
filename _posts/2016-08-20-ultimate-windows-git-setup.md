---
layout: post
title: "The Ultimate Git on Windows Setup"
date: 2016-08-20 13:37:00
comments: false
description: A guide for setting up Git in Windows 10 including PowerShell customization, multiple Git accounts and more.
modified: 2016-06-22
psversion: 5.1.14393.0
---

Source control and Git keeps getting more and more important for both Developers and Operations guys. Getting up and running with Git on MacOS or Linux is very easy as everything is built in.

On Windows, it's a bit of a different story. Let's spend a little time installing Git and customizing it to take our prompt from something that looks like this:

< pic of standard PS prompt >

to this:

< epic PS prompt >

* TOC
{:toc}

# Install Required Components

We will be installing the following tools for our ultimate git setup:

* `Chocolatey` - a Windows package manager
* Chocolatey Packages
  * `git.install` - Git for Windows
  * `ConEmu` - Terminal Emulator for Windows
* PowerShell Modules
  * `posh-git` - PowerShell functions for working with Git

Open an Administrative PowerShell prompt and enter the following:

{% highlight powershell %}
# Set your PowerShell execution policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force

# Install Chocolatey
iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex

# Install Chocolatey Packages
choco install git.install -y
choco install conemu -y

# Install PowerShell Modules
Install-PackageProvider NuGet -MinimumVersion '2.8.5.201' -Force
Set-PSRepository -Name PSGallery -InstallationPolicy Trusted
Install-Module -Name 'posh-git'
{% endhighlight %}

Close out of your PowerShell window.

# ConEmu

Open up ConEmu. I like to use this instead of the standard PowerShell prompt.

On the first launch of ConEmu, you will be prompted with a fast configuration dialog. Click `OK` and continue. We will customize it manually.

Open up the settings menu and configure the below settings.

![ConEmu Settings](/images/posts/windows_git/conemu_settings.png)

## Enable Single Instance Mode

Prevent multiple copies of ConEmu starting. Use the tabs instead!

![ConEmu Single Instance Mode](/images/posts/windows_git/conemu_settings_1.png)

## Enable Quake Style

This is a cool one, it makes ConEmu slide down from the top of your screen like the Quake terminal used to.

![ConEmu Quake Mode](/images/posts/windows_git/conemu_settings_2.png)

## Set PowerShell as the Default Shell

Who uses `cmd` anymore? Set the default shell to PowerShell.

![ConEmu PowerShell as default](/images/posts/windows_git/conemu_settings_3.png)

## Verify Quake Mode HotKey

Get the most out of Quake Mode by setting a hotkey.

![ConEmu PowerShell as default](/images/posts/windows_git/conemu_settings_4.png)

## Set a Custom Color Scheme

You can customize ConEmu you a color scheme. Check out the [ConEmu Theme GitHub Repo](https://github.com/joonro/ConEmu-Color-Themes). My terminal example above is using the `Dracula` theme.
