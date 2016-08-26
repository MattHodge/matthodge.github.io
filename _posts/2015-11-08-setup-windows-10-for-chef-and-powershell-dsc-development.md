---
layout: post
title:  "Setup Windows 10 For Chef and PowerShell DSC Development"
date:   2015-11-08 13:37:00
comments: true
modified: 2016-02-01
---

I am in the process of writing up some blog posts about working with PowerShell Desired State Configuration (DSC) and OpsCode Chef from a Windows Workstation / Windows Server perspective.

This first article will cover the steps required to setup a development environment for someone that is doing work with PowerShell Desired State Configuration (DSC) or OpsCode Chef.

We will be covering the following:

* Installing **git** and **poshgit** - source control for your code and PowerShell integration for git
* Installing **VirtualBox** - will be used with Vagrant to allow testing of DSC scripts and Chef recipes. I am using VirtualBox instead of Hyper-V as there are a ton of pre-build images for Linux and Windows available to download [Atlas](https://atlas.hashicorp.com/boxes/search). Not having to make these yourself is a huge time saver.
* Installing **Vagrant** - simple virtual machine creation. If you haven't used Vagrant before I recommend checking out [this video](https://www.youtube.com/watch?v=aUew6WauUsI).
* Installing **ChefDK** - the Chef Development Kit, includes knife, berkshelf, test-kitchen and Foodcritic
* Installing **Atom** - awesome text editor that includes linters which will check your code for issues as you type

Once everything is installed, we will be customizing and setting up the following:

* Customizing the **PowerShell Profile** to make it easier to work with git and set environment paths
* Customizing **Atom** with some plugins including a ruby linter to validate the ruby we write in the chef recipes
* Setting up **ssh keys** for use with git including configuring **ssh agent**, allowing us to push changes to git or work with private repositories
* Downloading some **vagrant plugins** and adding some vagrant boxes
* Configuring the `knife.rb` for working with Chef and creating a Berkshelf `config.json` file

Let's get started.

* TOC
{:toc}

## Set Execution Policy and Install Applications

The first step is to install the tools and applications we need. The easiest way to do this is with a combination of [Chocolatey](https://chocolatey.org/) and [OneGet](https://github.com/OneGet/oneget) which is now built into Windows 10.

* Open an Administrative PowerShell Prompt and do the following:

{% gist ff7d189b67a7b64c2ddf %}

## Customize Your PowerShell Profile

Now that you have installed [PoshGit](https://github.com/dahlbyk/posh-git), it is a good idea to take a look at your PowerShell Profile.

You can make any sort of customization you like, but here are two suggestions for things to do in your profile:

* Refresh the user and machine path environment variables. You will need to do this to use tools from the command line especially Atom, Chef and Vagrant.
* Perform a `Set-Location` on the directory you use the most – for me this is my `ProjectsGit` directory where I keep all of my git repositories.

I have posted [my PowerShell profile](https://github.com/MattHodge/MattHodgePowerShell/blob/master/PowerShellProfile/Microsoft.PowerShell_profile.ps1) on GitHub as an example [here](https://github.com/MattHodge/MattHodgePowerShell/blob/master/PowerShellProfile/Microsoft.PowerShell_profile.ps1).

My profile contains some useful functions for setting and reloading environment variables which I recommend you use in your own profile.

## Path Environment Variable

With many of the tools used for developing with Chef on Windows, they will require a correctly configured Path environment variable.

If one of your tools is not working or you cannot run it from the command line, there is a good chance something is wrong with your Path variable. For example, to use Ruby from the command line after you install the ChefDK, you need to [add ruby to the path variable](https://docs.chef.io/install_dk.html#add-ruby-to-path).

If you are using the functions from my [PowerShell profile](https://github.com/MattHodge/MattHodgePowerShell/blob/master/PowerShellProfile/Microsoft.PowerShell_profile.ps1) from above, it is very easy to add Path environment variables:

{% gist 4988ffcb68c8f4da44a4 %}

## Move PowerShell Profile to a Synced Drive (Optional)

If you move around to different machines, it is a good idea to move your your PowerShell Profile into a synced directory like Dropbox or OneDrive. From there you can create a symlink from the profile .ps1 in the synced path to the `$PROFILE` path.

* Move your PowerShell profile from the `$PROFILE` location to the synced folder
* Open an Administrative PowerShell Prompt
* Run the following commands:

{% gist d8ac66620dbe98717ea2 %}

![Symlink to PowerShell Profile](/images/posts/win10_for_chef_and_dsc/mklink_powershell_profile.png "Symlink to PowerShell Profile")

## Configure Atom

Atom is a great text editor which I like to use for everything but working on PowerShell scripts. It is my editor of choice for modifying Chef recipes.

We can make it more powerful with some additional packages.

* Open PowerShell
* Run the following commands to install some handy Atom packages

{% gist 35c585bfd3f59b4c6610 %}

## Setup Git SSH Keys

You should be using source control for your Chef recipes and PowerShell scripts. Warren Frame has an excellent blog on the topic specific to PowerShell [here](https://ramblingcookiemonster.github.io/GitHub-For-PowerShell-Projects/).

**Update 26/08/2016** - I created a more detailed guide on setting up Git SSH Keys on Windows. Check it out here: [Ultimate PowerShell Prompt Customization and Git Setup Guide](https://hodgkins.io/ultimate-powershell-prompt-and-git-setup)

## Download Vagrant Boxes

Vagrant is an excellent way to test your DSC scripts or Chef recipes.

I like to use 2 boxes for my Windows 2012 R2 and Ubuntu testing with Chef/DSC. There are also several plugins that make using Vagrant even nicer.

We will install some plugins and pre-load the Vagrant boxes:

{% gist dd9f4e708c77314bcc7b %}

## Configure Chef and Berkshelf

Next step is to get your chef `user.pem` file sorted out. Chef has a how-to guide for this [here](https://docs.chef.io/install_dk.html#manually-w-o-webui).

Once you have your `.pem` file, we will setup the `knife.rb` and the berkshelf `config.json`.

{% gist b8862934d2fe9b22cc8f %}

* Here is an [example knife.rb](https://github.com/MattHodge/MattHodgePowerShell/blob/master/Chef/knife_example.rb) that I use.
* Here is an example [berkshef .config](https://github.com/MattHodge/MattHodgePowerShell/blob/master/Chef/berkshelf_config_example.json) file that I use.

## Customize the PowerShell ISE Theme (Optional)

The default theme for the PowerShell ISE is boring, lets spice it up with a theme. There is a great repository with PowerShell ISE themes located [on GitHub](https://github.com/marzme/PowerShell_ISE_Themes).

To import the themes into the PowerShell ISE, go to **Tools > Options > Manage Themes > Import**.

![PowerShell ISE Themes](/images/posts/win10_for_chef_and_dsc/powershell_ise_themes.png "PowerShell ISE Themes")

Again, I would drop the theme into a synced folder so you can use it on all your machines.

## Conclusion

With that, you should have your Windows machine setup PowerShell DSC and Chef development. Did I miss anything? Send me a tweet [@matthodge](https://twitter.com/matthodge) and let me know!
