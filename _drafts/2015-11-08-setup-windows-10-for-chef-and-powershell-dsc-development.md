---
layout: post
title:  "Setup Windows 10 For Chef and PowerShell DSC Development"
date:   2015-11-08 13:37:00
comments: false
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

{% highlight powershell %}
# Configure PowerShell Execution Policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force

# Install Chocolatey
iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1'))

# Install the required apps
choco install git.install -y
choco install virtualbox -y
choco install vagrant -y
choco install chefdk -y
choco install atom -y
choco install poshgit -y

# Optional - install tabbed Explorer
choco install clover -y

# Optional - free git GUI
choco install sourcetree -y

# Update Pester - https://github.com/pester/Pester
Install-Module -Name Pester -Confirm:$true

# Install the PSScriptAnalyzer - https://github.com/PowerShell/PSScriptAnalyzer
Install-Module -Name PSScriptAnalyzer -Confirm:$true
{% endhighlight%}

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

If you are using the functions from my [PowerShell profile](https://github.com/MattHodge/MattHodgePowerShell/blob/master/PowerShellProfile/Microsoft.PowerShell_profile.ps1)from above, it is very easy to add Path environment variables:

{% highlight powershell %}
# Requires the Add-PathVariable function from my PowerShell Profile
Add-PathVariable -Path 'C:/opscode/chefdk/embedded/bin'
{% endhighlight%}

## Move PowerShell Profile to a Synced Drive (Optional)

If you move around to different machines, it is a good idea to move your your PowerShell Profile into a synced directory like Dropbox or OneDrive. From there you can create a symlink from the profile .ps1 in the synced path to the `$PROFILE` path.

* Move your PowerShell profile from the `$PROFILE` location to the synced folder
* Open an Administrative PowerShell Prompt
* Run the following commands:

{% highlight powershell %}
# Create a symlink to the profile in your shared drive
cmd /c mklink $PROFILE D:\DataHodge\Dropbox\PSProfile\Microsoft.PowerShell_profile.ps1

# Load the profile into the current session
. $PROFILE
{% endhighlight%}

![Symlink to PowerShell Profile](/images/posts/win10_for_chef_and_dsc/mklink_powershell_profile.png "Symlink to PowerShell Profile")

## Configure Atom

Atom is a great text editor which I like to use for everything but working on PowerShell scripts. It is my editor of choice for modifying Chef recipes.

We can make it more powerful with some additional packages.

* Open PowerShell
* Run the following commands to install some handy Atom packages

{% highlight powershell %}
# Linter to validate the code as you are typing
apm install linter

# Install rubocop gem
gem install rubocop

# Linter for ruby
apm install linter-rubocop

# Rubocop auto corrector
apm install rubocop-auto-correct

# Create a rubocop.yml configuration file to ignore warnings for line endings. Details here https://github.com/bbatsov/rubocop/blob/master/README.md
Set-Content -Path ~/.rubocop.yml -Value 'Metrics/LineLength:','  Enabled: false'

# Useful for removing Windows line endings
apm install line-ending-converter

# Gives a view of your entire document when it is open in atom
apm install minimap

# monokai theme for atom
apm install monokai
{% endhighlight%}

## Setup Git SSH Keys

You should be using source control for your Chef recipes and PowerShell scripts. Warren Frame has an excellent blog on the topic specific to PowerShell [here](https://ramblingcookiemonster.github.io/GitHub-For-PowerShell-Projects/).

To work with git repositories, it is best to use ssh keys. On Windows, the ssh keys live under your user directory in a `.ssh` folder, for example `C:\Users\YourName\.ssh`

Setting up ssh keys on Windows for GitHub and BitBucket can be a bit of a pain, but the below will guide you through the process

## Download Vagrant Boxes

Vagrant is an excellent way to test your DSC scripts or Chef recipes.

I like to use 2 boxes for my Windows 2012 R2 and Ubuntu testing with Chef/DSC. There are also several plugins that make using Vagrant even nicer.

We will install some plugins and pre-load the Vagrant boxes:

{% highlight powershell %}
# Install vagrant plugins
vagrant plugin install 'vagrant-berkshelf'
vagrant plugin install 'vagrant-dsc'
vagrant plugin install 'vagrant-omnibus'
vagrant plugin install 'vagrant-reload'
vagrant plugin install 'vagrant-vbguest'
vagrant plugin install 'vagrant-vbox-snapshot'
vagrant plugin install 'vagrant-winrm'
vagrant plugin install 'winrm-fs'

# Install vagrant boxes
vagrant box add ubuntu/trusty64
vagrant box add kensykora/windows_2012_r2_standard

# Install the test-kitchen plugins
gem install kitchen-pester
{% endhighlight%}

## Configure Chef and Berkshelf

Next step is to get your chef `user.pem` file sorted out. Chef has a how-to guide for this [here](https://docs.chef.io/install_dk.html#manually-w-o-webui).

Once you have your `.pem` file, we will setup the `knife.rb` and the berkshelf `config.json`

* Here is an [example knife.rb](https://github.com/MattHodge/MattHodgePowerShell/blob/master/Chef/knife_example.rb) that I use.

* Here is an example [berkshef .config](https://github.com/MattHodge/MattHodgePowerShell/blob/master/Chef/berkshelf_config_example.json) file that I use.

## Customize the PowerShell ISE Theme (Optional)

The default theme for the PowerShell ISE is boring, lets spice it up with a theme. There is a great repository with PowerShell ISE themes located [on GitHub](https://github.com/marzme/PowerShell_ISE_Themes).

To import the themes into the PowerShell ISE, go to **Tools | Options | Manage Themes | Import**.

<a href="http://www.hodgkins.net.au/wp-content/uploads/2015/11/image.png" rel="lightbox[753]"><img style="background-image: none; float: none; padding-top: 0px; padding-left: 0px; margin: 0px auto; display: block; padding-right: 0px; border: 0px;" title="image" src="http://www.hodgkins.net.au/wp-content/uploads/2015/11/image_thumb.png" alt="image" width="945" height="591" border="0" /></a>

Again, I would drop the theme into a synced folder so you can use it on all your machines.

## Conclusion

With that, you should have your Windows machine setup PowerShell DSC and Chef development. Did I miss anything? Send me a tweet <a href="https://twitter.com/matthodge" target="_blank">@matthodge</a> and let me know!
