---
layout: post
title: "Ultimate PowerShell Prompt and Git Setup"
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

# Install required components

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

# Install Chocolatey packages
choco install git.install -y
choco install conemu -y

# Install PowerShell modules
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

## Enable single instance Mode

Prevent multiple copies of ConEmu starting. Use the tabs instead!

![ConEmu Single Instance Mode](/images/posts/windows_git/conemu_settings_1.png)

## Enable Quake mode

This is a cool one, it makes ConEmu slide down from the top of your screen like the Quake terminal used to.

![ConEmu Quake Mode](/images/posts/windows_git/conemu_settings_2.png)

## Set PowerShell as the default shell

Who uses `cmd` anymore? Set the default shell to PowerShell.

![ConEmu PowerShell as default](/images/posts/windows_git/conemu_settings_3.png)

## Verify Quake mode hot-key

Get the most out of Quake Mode by setting a hotkey.

![ConEmu PowerShell as default](/images/posts/windows_git/conemu_settings_4.png)

## Set a custom color scheme

You can customize ConEmu you a color scheme. Check out the [ConEmu Theme GitHub Repo](https://github.com/joonro/ConEmu-Color-Themes). My terminal example above is using the `Dracula` theme.

# PowerShell Profile

We have a nice terminal theme, but let's do a few finishing touches to make it pop.

## Create and edit the PowerShell Profile

PowerShell can load some settings every time it starts, which is known as the PowerShell Profile or `$PROFILE`.

To create/edit your `$PROFILE` do the following:

{% highlight powershell %}
# Creates profile if doesn't exist then edits it
if (!(Test-Path -Path $PROFILE)){ New-Item -Path $PROFILE -ItemType File } ; ise $PROFILE
{% endhighlight %}

This will launch the PowerShell ISE so you can edit the profile.

## Import the posh-git Module
The first thing to do inside your PowerShell Profile is to import the `posh-git` module.

Add the following to your `$PROFILE`

{% highlight powershell %}
Import-Module -Name posh-git
{% endhighlight %}

This will customize our prompt when we are inside git repos.

## Customize the prompt

Let's make our prompt a little cooler and customize it a little.

I like the prompt that `Joon Ro` created over at [his blog](https://joonro.github.io/blog/posts/powershell-customizations.html). I modified it slightly:

{% gist 0f0be96e0489feeb8a05d94151093517 %}

## Colorize your directory listing

When we do a `ls` or `dir` wouldn't it be nice to be able to colorize folders or certain file types instead of just having a boring list that looks the same?

Check out the [Get-ChildItem-Color](https://github.com/joonro/Get-ChildItem-Color) repository. I added the contents of `Get-ChildItem-Color.ps1` to my `$PROFILE`.

I then overwrote both the `ls` and `dir` aliases by adding the following into my `$PROFILE`:

{% highlight powershell %}
Set-Alias ls Get-ChildItem-Color -option AllScope -Force
Set-Alias dir Get-ChildItem-Color -option AllScope -Force
{% endhighlight %}

# Git

Now we have a nice terminal to work with, let's get Git setup.

Open up ConEmu.

## Add C:\Program Files\Git\usr\bin to Path Variable

First up we need to add the `C:\Program Files\Git\usr\bin` folder to our path variable. This folder contains `ssh-add` and `ssh-agent` which we will be using to manage our SSH keys.

{% highlight powershell %}
# Permanently add C:\Program Files\Git\usr\bin to machine Path variable
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\Git\usr\bin", "Machine")
{% endhighlight %}

Restart ConEmu for it to take effect.

## Generate a key

Let's generate our ssh key.

{% highlight powershell %}
# Generate the key and put into the your user profile .ssh directory
ssh-keygen -t rsa -b 4096 -C "your@email.com" -f $env:USERPROFILE\.ssh\id_rsa
{% endhighlight %}

## Add the public key to GitHub

Once we have a generated SSH Key, we need to give GitHub the public key.

{% highlight powershell %}
# Copy the public key. Be sure to copy the .pub for the public key
Get-Content $env:USERPROFILE\.ssh\id_rsa.pub | clip
{% endhighlight %}

Open up your GitHub settings and choose `SSH and GPG keys` on the left.

![Add Github Public Key](/images/posts/windows_git/github_add_new_key.png)

This process is similar for BitBucket.

## Add our key to ssh-agent

When we try and push to our git repository, our machine will need to authenticate us using our SSH Key. A tool called `ssh-agent` keeps track of the keys we have and authenticating against GitHub for us.

# More Advanced Git Usage
