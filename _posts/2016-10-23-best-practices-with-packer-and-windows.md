---
layout: post
title: Best Practices with Packer and Windows
date: 2016-10-25T13:37:00.000Z
comments: true
description: Some best practices around building Windows based Packer images.
---

* TOC
{:toc}

# Why you should be using Packer

Already know why Packer is useful? Jump directly to the [best practices](#best-practices).

When you develop automation, for example, PowerShell Desired State configuration resources, where do you test them?

If the answer is locally on your machine or a remote Virtual Machine platform, you are missing out on some opportunities of speed and reduction in your development and test cycle time.

Have a look at [Gael Colas's](https://twitter.com/gaelcolas) awesome introduction to [Test-Kitchen and Kitchen-DSC](https://gaelcolas.com/2016/07/11/introduction-to-kitchen-dsc/), which will show you how to develop and test your DSC resources easily on your local machine inside Virtual Machines.

As part of this workflow, you will need to use a base virtual machine where you apply your DSC configurations too. Depending on your environment, you may need to apply the same DSC resource to multiple machine configurations, for example:

* Windows 2012 R2 WMF 4.0 with GUI
* Windows 2012 R2 WMF 4.0 Core
* Windows 2012 R2 WMF 5.0 with GUI
* Windows 2012 R2 WMF 5.0 Core
* Windows 2008 R2 WMF 4.0 with GUI
* Windows 2008 R2 WMF 4.0 Core
* Windows 2008 R2 WMF 5.0 with GUI
* Windows 2008 R2 WMF 5.0 Core
* Windows 2016 WMF 5.0 with GUI
* Windows 2016 WMF 5.0 Core

That is 10 different variations of Windows you need to maintain templates for! Do you manually go through each of them every month to apply their Windows updates too? What about if you need to share these base images with colleges? Do you copy 50GB of images over the internet or make your colleges build their own images? What if one of your colleges uses a Mac and has to use VirtualBox images instead of Hyper-V?

You started using Test-Kitchen because it was meant to simplify your workflow and now you have an image management problem!

![Packer Logo](/images/posts/packer_best_practices/packer_logo.png)

This is where [Packer](https://www.packer.io/) by [HashiCorp](https://www.hashicorp.com/) can help. Packer is a tool for creating machine images from a single configuration source. You store the entire image creation process as code, so images are always built the same way, this way, instead of having to ship entire VM templates over the internet, you can just keep your Packer configuration in source control and anyone in your team can build their own templates.

# Getting Started with Packer

I will assume some base knowledge of Packer for this guide. If you are just getting started with it, I recommend [Matt Wrock's](https://twitter.com/mwrockx/) blog post titled [Creating windows base images using Packer and Boxstarter](http://www.hurryupandwait.io/blog/creating-windows-base-images-for-virtualbox-and-hyper-v-using-packer-boxstarter-and-vagrant).


# Best Practices

## Step by Step

:white_check_mark: **When creating Packer templates, create builds in a step by step process. Do not try and do everything in a single Packer build.**

When I first started creating my Packer templates for Windows, I would include everything in a single `.json` file:

* Build the Windows box from a `.iso`
* Apply Windows Updates
* Install Windows Management Framework 5.0
* Convert to a Vagrant Box
* Upload to Atlas

The problem with this is if at any point there is a failure, you need to start the whole build process again, as Packer will automatically stop and delete a Virtual Machine on failure.

This is particularly annoying when you have just waited 4+ hours for Windows Updates to occur. There were more than a few times I felt like this having things fail after Windows updates completed:

![Laptop Smash](https://i.imgur.com/UJqBiL6.gif)

Instead, create several Packer build `.json` files and "chain" them together:

1. From Windows `.iso` to a working machine using an `Autounattand.xml` answer file
1. From base image to Windows Updates and WMF 5.0
1. From updated image to cleaned image (Defrag, Remove temp files etc.)

Doing this gives you several benefits:

* If a step fails, you can resume from a previous step, not start the entire process again which will save you lots of time
* It allows you to create some branching logic from your builds, for example:

![Packer Branch Build](/images/posts/packer_best_practices/packer_branching.png)

In this branching example, the matching colors mean that the same template is used. This would mean you could end up with 2 base images, one with no updates, and one with updates and WMF 5.0 installed.

The way this works is to have the Packer builder output an image with a specified name to a specified directory.

Here is an example using the `virtualbox-iso` and `virtualbox-ovf` builders. The first build in the chain coming from an .`iso` file:

{% highlight json %}
{
  "builders": [
    {
      "type": "virtualbox-iso",
      "output_directory": "./output-window2012r2-base/",
      "vm_name": "win2012r2-base"
    }
  ]
}
{% endhighlight %}

The next build in the chain coming from the previous build's VirtualBox `.ovf` file.

{% highlight json %}
{
  "builders": [
    {
      "type": "virtualbox-ovf",
      "source_path": "./output-window2012r2-base/win2012r2-base.ovf",
      "output_directory": "./output-window2012r2-with-updates/",
      "vm_name": "win2012r2-with-updates"
    }
  ],
  "provisioners": [
    {
      "type": "powershell",
      "inline": [
        "Do windows updates here"
      ]
    }
  ]
}
{% endhighlight %}

You can read in detail about [User Variables](https://www.packer.io/docs/templates/user-variables.html) in the [Packer Docs](https://www.packer.io/docs).

## Generic Templates

Once you start making a few Packer `.json` template files for different Windows versions, you will notice they start to become very similar, and you will find you are repeating a lot of your code.

:white_check_mark: **Keep templates as generic as possible, and use User Variables.**

When your templates are generic and accept user variables, you can pass variables to the Packer template via the command line. In the below example, I am passing the source path for the base image (`C:\packer\output-window2012r2-basewin2012r2-base.ovf`) into the build.

{% gist 98e402742eb8eb739ce8d624092405a3 %}

Here is a snippet of `02-win_updates-wmf5.json`:

{% highlight json %}
{
  "builders": [
    {
      "type": "virtualbox-ovf",
      {%- raw -%}
      "source_path": "{{user `source_path`}}"
      {% endraw -%}
    }
  ],
  "variables": {
    "source_path": ""
  }
}
{% endhighlight %}

Doing this also makes it possible to reuse the same Packer build template at any part of your build branch (see above).

## Use guest additions mode of attach

:white_check_mark: **When installing guest additions, use attach mode.**

Using the `attach` method for the guest additions is much faster and more reliable than using upload, especially over WinRM.

You will also need to un-mount the attached ISO after the build completes using `vboxmanage_post` because the ISO will not detach itself. This issue is being tracked on GitHub in issue [#3121](https://github.com/mitchellh/packer/issues/3121).

{% highlight json %}
{
  "builders": [
    {
      "guest_additions_mode": "attach",
      "vboxmanage_post": [
        [
          "storageattach",
          {%- raw -%}
          "{{.Name}} }}",
          {% endraw -%}
          "--storagectl",
          "IDE Controller",
          "--port",
          "1",
          "--device",
          "0",
          "--medium",
          "none"
        ]
      ]
    }
  ],
  "provisioners": [
    {
      "type": "powershell",
      "script": "scripts/install_oracle_guest_additions.ps1",
      "elevated_user": "vagrant",
      "elevated_password": "vagrant"
    }
  ]
}
{% endhighlight %}

Then you can use the PowerShell script `install_oracle_guest_additions.ps1` to install the tools:

{% gist 0f0c159d68fdf80d6e66659bcfbac83a %}

## Use environment variables to change the action of a provisioning script

This one helps you to keep your templates very generic. Say for example I wanted to have the option to install VirtualBox tools or not for a build, without having separate templates, I could do the following:

* Create a user variable called `install_vbox_tools` which accepts either true or false
* Pass the user variable to the provisioning script where I am running `install_oracle_guest_additions.ps1` as an environment variable called `install_vbox_tools`
* Inside the `install_oracle_guest_additions.ps1`, read the `install_vbox_tools` environment variable and take either install them or not depending on if the variable is true or false

Here is an example Packer template:

{% gist 61f3ecb0a755afdc5ff4bc18c5af990a %}

This is an example of the `install_oracle_guest_additions.ps1` PowerShell script:

{% gist 9e7d50714bda3ca6b9f02fe0f5b9e166 %}

When you run the Packer build, you can pass `-var "install_vbox_tools=true"` or `-var "install_vbox_tools=false"` and the PowerShell script provisioner will take the appropriate action.

## Keep the OS information in your build script

:white_check_mark: **Keep your operating system related information inside your build scripts instead of inside the Packer templates.**

Now you have made your Packer build `.json` files very generic, you can move the OS related information into a build a script, and pass them as user variables.

Here is an example PowerShell build script, where the Windows 2012 R2 or Windows 2016 Core could be installed using the same Packer templates.

{% gist cc1214f192696c5754975c2c75b7a5e7 %}

## Disable WinRM on build completion and only enable it on first boot

If you are running `sysprep` on your Windows images, when they first boot they will need to restart themselves. This fine in normal circumstances, but when using the images in Vagrant for example, on the first boot, Vagrant will detect that WinRM is up and start connecting, and then the machine will restart. This will make Vagrant think the machine has failed or isn't in the correct state.

The trick to this is having WinRM disabled until the very last moment, after the initial sysprep reboot.

:white_check_mark: **Keep WinRM disabled or blocked by the firewall until the system has had its final boot after sysprep.**

To do this as part of your Packer build:

* Use a PowerShell script provisioner to drop a `PackerShutdown.bat` file on the system. This shutdown command will block WinRM in the firewall and then sysprep the machine. We will use this as the Packer shutdown command from inside our build.

{% gist 3c738a165a1b9afc1b3249cf937a53a8 %}

* Use a PowerShell script provisioner to create a batch file at `C:\Windows\Setup\Scripts\SetupComplete.cmd` containing the following:

{% gist 737406f6d124006019e0f6d1d75e5dfe %}

This script will be run automatically by Windows after the machine reboots from the sysprep. It will unblock WinRM in the firewall at the right time for us.

* Use the `PackerShutdown.bat` file we created as the `shutdown_command` in your builder:

{% highlight json %}
{
  "builders": [
    {
      "shutdown_command": "C:/Windows/packer/PackerShutdown.bat",
      "shutdown_timeout": "1h"
    }
  ],
  "provisioners": [
    {
      "type": "powershell",
      "script": "scripts/save_shutdown_command.ps1"
    }
  ]
}
{% endhighlight %}

You can find a full example of the `save_shutdown_command.ps1` script in [my Packer repo](https://github.com/MattHodge/PackerTemplates).

## Use headless mode

:white_check_mark: **Use headless mode when building images.**

When using the VirtualBox builder, using headless mode errors a lot less. You can still access the GUI for the virtual machine by manually loading the VirtualBox and double clicking on it.

{% highlight json %}
{
  "builders": [
    {
      "type": "virtualbox-iso",
      "headless": "true"
    }
  ]
}
{% endhighlight %}


## Set a high shutdown and WinRM timeouts

As you have probably noticed when you install a ton of Windows updates on a machine, it can take a long time to reboot.

To prevent this sort of thing from causing the Packer build to fail, make sure you set a high `winrm_timeout` and `shutdown_timeout`.

{% highlight json %}
{
  "builders": [
    {
      "winrm_timeout": "12h",
      "shutdown_timeout": "1h"
    }
  ]
}
{% endhighlight %}

# Conclusion

Using these practices will help you on your way to creating some awesome Windows based Packer images.

All the code mentioned in this post is available at [https://github.com/MattHodge/PackerTemplates/](https://github.com/MattHodge/PackerTemplates/).

Thanks to [@mwrockx](https://twitter.com/mwrockx/) for his work getting me started on Packer templates and [@kikitux](https://twitter.com/kikitux) for his advice on Packer.
