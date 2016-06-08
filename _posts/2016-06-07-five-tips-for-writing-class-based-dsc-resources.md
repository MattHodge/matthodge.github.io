---
layout: post
title: "5 Tips for Writing DSC Resources in PowerShell v5"
date: 2016-06-08 13:37:00
comments: false
description: Five useful tips when writing DSC Resources using PowerShell v5. Covers folder structure, IntelliSense, verifying exposed resources and testing using Pester.
psversion: 5.0.10586.122
---

I have recently been writing some class based DSC resources and have been enjoying the experience much more than (trying to) write DSC resources in PowerShell v4. Using classes really simplifies the development process of the DSC resources, and I believe class based resources will be the norm going forward.

As PowerShell v5 RTM has only been recently released back in February 2016, the large majority of DSC resources available on GitHub still use the `Get-TargetResource`, `Test-TargetResource` and `Set-TargetResource` functions, so I found it a little hard to get started as there was a lack of examples.

Hopefully the following tips will save you some time and pain when going down the path of writing your own class based resources!


* TOC
{:toc}

## Tip 1 - Structure your DSC Resource in a standard manner

When working on anything PowerShell related, you should be using source control. If you aren't, you had better take a [crash course in version control](https://www.youtube.com/watch?v=wmPfDbsPeZY) with  Warren Frame  ([@pscookiemonster](https://twitter.com/psCookieMonster)).

I personally have a `C:\ProjectsGit` directory that I create all my repos inside it.

This is how I structure my DSC resources, for an example resource called `MyDSCResource` would be structured like this:

```
C:\ProjectsGit
└── MyDSCResource\                       # Repo Directory
    ├── MyDSCResource\                   # DSC Resource Directory
    |   ├── MyDSCResource.psd1           # Manifest File
    |   └── MyDSCResource.psm1           # Class Based Resource
    ├── Examples\                        # Resource usage examples
    |   └── Example1.ps1
    ├── Tests\                           # Pester Tests
    ├── appveyor.yml OR .build.ps1       # Build file
    └── readme.md
```

If you need a composite resource to go with your DSC Resource, the structure would be as follows:

```
C:\ProjectsGit
└── MyDSCResource2\
    ├── MyDSCResource2\
    |   ├── DSCResources\                  # Mandatory folder name for composite resources
    |   |   └── CompResourceName\          # Name of composite resource
    |   |       ├── CompResourceName.psd1  # Manifest File for composite resource
    |   |       └── CompResourceName.psm1  # Composite Resource      
    |   ├── MyDSCResource2.psd1
    |   └── MyDSCResource2.psm1
    ├── Examples\
    |   └── Example1.ps1
    ├── Tests\
    ├── appveyor.yml OR .build.ps1
    └── readme.md
```

:white_check_mark: **Mini Tip:** You cannot write composite resources as a class-based resource, but you can include composite resources WITH class based resources, take this example:

![Class based DSC Resource with composite resource](/images/posts/five_dsc_tips/composite_resource_class_based.png)

## Tip 2 - Create a symbolic link between the DSC resource in your local repo and your PowerShell module path

When writing DSC configurations against your custom resources in the PowerShell IDE, it makes it much easier to have IntelliSense working on your resources. The way this is usually achieved is by copying your custom resource into one of the `$env:PSModulePath` folder; for example the `C:\Program Files\WindowsPowerShell\Modules` folder.

If you don't have your custom resource in your module folder, you get this:

![Resource not in Module Path](/images/posts/five_dsc_tips/resource_not_in_module_path.png)

Since we want to have our resource in source control at all times, it can become painful to remember to copy your resource into the modules path every time every time you make a change.

A simple solution to this is to make a symbolic link between your resource repository and the desired path in the `Modules` folder:

{% highlight powershell %}
# originalPath is the one containing the .psm1 and .psd1
$originalPath = 'C:\ProjectsGit\MyDSCResource1\MyDSCResource1\'

# pathInModuleDir is the path where the symbolic link will be created which points to your repo
$pathInModuleDir = 'C:\Program Files\WindowsPowerShell\Modules\MyDSCResource1'

New-Item -ItemType SymbolicLink -Path $pathInModuleDir -Target $originalPath
{% endhighlight %}

:white_check_mark: **Mini Tip:** The ISE will sometimes not automatically realize your resource has appeared in the modules directory. To fix this, just delete and replace the `Import-DSCResource -ModuleName MyDSCResource1` line in the ISE and it will discover your module correctly and give you IntelliSense.

## Tip 3 - Press Control + Space For IntelliSense on the DSC Resource

When writing DSC configurations, you can use IntelliSense against DSC resource's (including your custom ones) by pressing `Control + Space` inside the DSC configuration block for the resource.

![Use Control + Space in ISE](/images/posts/five_dsc_tips/press_control_space.png)

## Tip 4 - Verifying the resources your DSC module is exposing

An easy way to verify the DSC resources your custom DSC module is exposing is to use the  `Get-DSCResource` CmdLet:

![Do a Get-DSCResource](/images/posts/five_dsc_tips/use_get_dscresource.png)

This allows you to confirm that your resource is working correctly.

If for some reason your resource is not being exposed, verify if you added it to the `DscResourcesToExport` array inside the DSC module manifest file (`.psd1`):

{% highlight powershell %}
DscResourcesToExport = @('HubotInstall','HubotInstallService')
{% endhighlight %}

If you are trying to create a composite resource, part of the process is to create a module manifest file for the composite. If for some reason your composite resource is not appearing in the list when you perform a `Get-DSCResource`, try manually loading the manifest for the composite resource.

As an example, say I have my DSC Resource and composite structured like this:

```
C:\ProjectsGit
└── Hubot-DSC-Resource\
    ├── Hubot\
    |   ├── DSCResources\
    |   |   └── HubotPrerequisites\
    |   |       ├── HubotPrerequisites.psd1
    |   |       └── HubotPrerequisites.psm1
    |   ├── Hubot.psd1
    |   └── Hubot.psm1
    └── readme.md
```

I would use `Import-Module` to try and import the composite resource `HubotPrerequisites.psd1`:

{% highlight powershell %}
 Import-Module C:\ProjectsGit\Hubot\Hubot\DSCResources\HubotPrerequisites\HubotPrerequisites.psd1
{% endhighlight %}

This will allow me to confirm if the module is loading correctly, or give a reason why it cannot be loaded. Getting an import error means your module won't expose the composite resource:

![DSC Composite Resource Troubleshooting](/images/posts/five_dsc_tips/error_when_importing_composite_resource.png)

## Tip 5 - Testing Class Based Resources with Pester

You are able to test your class based resources by putting the following in your Pester tests:

{% highlight powershell %}
# 'using module' on the class based DSC resource .psm1
using module ..\Hubot\Hubot.psm1
{% endhighlight %}

Putting the `using module` is the secret sauce to allow you to access your classes and the properties and methods inside them.

Here is an example of how to access the test block of my custom `[HubotInstall]` resource:

{% highlight powershell %}
# using on your DSC Resource .psm1
using module ..\Hubot\Hubot.psm1

# load a class in my DSC Resource
$x = [HubotInstall]::new()

# set the parameters required by the class to do a set/test/get
$x.BotPath = 'C:\mybot'
$x.Ensure = 'Present'

# Execute the test method of the [HubotInstall] class
$x.Test()
{% endhighlight %}

Take a look at [Chris Hunt's](https://twitter.com/logicaldiagram) blog post [Testing PowerShell Classes](https://www.automatedops.com/blog/2016/01/28/testing-powershell-classes/) for more details on `using module`.

Unfortunately, with the current version of PowerShell, when you load a class using the `using module` method, it is cached into the PowerShell session. Even if you make changes to your class and re-run `using module`, the old class will be used. Take a look at this example:

![PowerShell class caching](https://i.imgur.com/Q10DMf6.gif "PowerShell Class Caching")

To get around this, we can use a function to make Pester invoke the tests inside its own run space, meaning that the class is loaded fresh each time. (Thanks to [Dave Wyatt](https://twitter.com/msh_dave) for this tip.)

{% gist 221b2478069f56f1bf95fe98e50a095c %}

Running your tests using `Invoke-PesterJob` provides a workaround for the caching issue.

## Wrapping up and further reading

Remember the 5 tips above when on your DSC resource writing journey, they could save you some time and pain when troubleshooting!

If you want to dive deeper into writing classed based resources, or learn about how to debug DSC Resources at run time in PowerShell v5, I recommend you check out [What's New in PowerShell v5](https://mva.microsoft.com/en-US/training-courses/whats-new-in-powershell-v5-16434) on Microsoft Virtual Academy.
