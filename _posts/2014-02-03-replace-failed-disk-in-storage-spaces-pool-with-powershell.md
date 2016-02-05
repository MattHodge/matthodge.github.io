---
layout: post
title:  "Replacing a Failed Disk in Windows Server 2012 R2 Storage Spaces with PowerShell"
date:   2014-01-06 13:37:00
comments: false
modified: 2016-02-05
---

Failed hard disks are in-evadable. There are many ways to provide resiliency for hard disk failure, and Windows Server 2012/Windows Server 2012 R2’s build in feature to provide this is Storage Spaces.

A hard disk failed inside my Storage Pool, so lets switch over to PowerShell to get this resolved.

* TOC
{:toc}

## Diagnosis

Firstly, open up an Administrative PowerShell prompt. To get the status of my Storage Space (which I called `pool`) I run the command:

{% highlight powershell %}
Get-StoragePool
{% endhighlight %}

![Get-StoragePool](/images/posts/storage_spaces_powershell/01_get_storagepool.png "Get-StoragePool")

I can see that my Storage Space named `pool` is in a degraded state.

To check the health of the volumes sitting inside the Storage Pool, use the command

{% highlight powershell %}
Get-VirtualDisk
{% endhighlight %}

![Get-VirtualDisk](/images/posts/storage_spaces_powershell/02_get_virtualdisk.png "Get-VirtualDisk")

We can see that `Media`, `Software` and `DocumentsPhotos` volumes are have `Degraded` as their `OperationalStatus`. This means that they are still attached and accessible, but their reliability cannot be ensured should there be another drive failure. These volumes have either a Parity or Mirror parity setting, which has allowed Storage Spaces to save my data even with the drive failure.

The `Backups` and `VMTemplates` have a `Detached` operational status. I was not using any resiliency mode on this data as it is easily replaced, so it looks like I have lost the data on these volumes.

To get an idea what is happening at the physical disk layer I run the command:

{% highlight powershell %}
Get-PhysicalDisk
{% endhighlight %}

![Get-PhysicalDisk](/images/posts/storage_spaces_powershell/03_get_physicaldisk.png "Get-PhysicalDisk")

We can see that `PhysicalDisk1` is in a failed state. As the HP N40L has a 4 bay enclosure with 4TB Hard Disks in them, it is easy to determine that PhyisicalDisk1 is in the first bay in the enclosure.

## Retiring the Failed Disk

Now I determined which disk had failed, the server was shutdown and the failed disk from the first bay was replaced with a spare 4TB Hard Disk.

With the server back online, open PowerShell back up with administrative permissions and check what the physical disks look like now:

{% highlight powershell %}
Get-PhysicalDisk
{% endhighlight %}

![Get-PhysicalDisk with drive replaced](/images/posts/storage_spaces_powershell/04_retired_failed_disk.png "Get-PhysicalDisk with drive replaced")

We can see that the new disk that was installed has taken the `FriendlyName` of `PhysicalDisk1` and has a `HealthStatus` of `Healthy`. The failed disk has lost its `FriendlyName` and its `OperationalStatus` has changed to `Lost Communication`.

First lets single out the missing disk:

{% highlight powershell %}
Get-PhysicalDisk | Where-Object { $_.OperationalStatus -eq 'Lost Communication' }
{% endhighlight %}

![Get-PhysicalDisk filtering out failed disk](/images/posts/storage_spaces_powershell/05_filter_failed_disk.png "Get-PhysicalDisk filtering out failed disk")

Assign the missing disk to a variable:

{% highlight powershell %}
$missingDisk = Get-PhysicalDisk | Where-Object { $_.OperationalStatus -eq 'Lost Communication' }
{% endhighlight %}

Next we need to tell the storage pool that the disk has been retired:

{% highlight powershell %}
$missingDisk | Set-PhysicalDisk -Usage Retired
{% endhighlight %}

## Adding a New Disk

To add the replacement disk into the Storage Pool

{% highlight powershell %}
# Set the new Disk to a Variable
$replacementDisk = Get-PhysicalDisk –FriendlyName PhysicalDisk1

# Add the disk to the Storage Pool
Add-PhysicalDisk –PhysicalDisks $replacementDisk –StoragePoolFriendlyName pool
{% endhighlight %}

## Repairing the Volumes
The next step after adding the new disk to the Storage Pool is to repair each of the Virtual Disks residing on it.

{% highlight powershell %}
# To Repair and Individual Volume
Repair-VirtualDisk –FriendlyName DocumentsPhotos

# To Repair all Warning Volumes
Get-VirtualDisk | Where-Object –FilterScript {$_.HealthStatus –Eq 'Warning'} | Repair-VirtualDisk
{% endhighlight %}

We can see the the repair running by entering:

{% highlight powershell %}
Get-VirtualDisk
{% endhighlight %}

![See disk repairing](/images/posts/storage_spaces_powershell/06_disk_repairing.png "See disk repairing")

The `OperationalStatus` of `InService` lets us know the volume is currently being repaired. The percentage completion of the repair can be found by running:

{% highlight powershell %}
Get-StorageJob
{% endhighlight %}

![Get-StorageJob](/images/posts/storage_spaces_powershell/07_get_storagejob.png "Get-StorageJob")

## Remove the Lost VirtualDisks

Since there were no parity on the `VMTemplates` and `Backups` Volumes, they can be deleted with the following command:

{% highlight powershell %}
Remove-VirtualDisk –FriendlyName <FriendlyName>
{% endhighlight %}

## Removing the Failed Disk from the Pool

This step will not work if you still have `Degraded` disk in the Storage Pool, so make sure all repairs complete first.

{% highlight powershell %}
Remove-PhysicalDisk –PhysicalDisks $missingDisk –StoragePoolFriendlyName pool
{% endhighlight %}

## Summary

To wrap up, to replace a failed disk in a storage pool:

{% gist 0d3401a52f2a7080719a %}

The full list of Windows Server Storage Spaces CmdLets can be found on TechNet here: [http://technet.microsoft.com/en-us/library/hh848705.aspx](http://technet.microsoft.com/en-us/library/hh848705.aspx).
