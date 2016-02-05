---
layout: post
title:  "How to fix a Timeout Erorr while mounting a LUN in NetApp SnapDrive"
date:   2012-05-30 13:37:00
comments: true
modified: 2016-02-05
---

When you are trying to mount a LUN (by using Connect Disk) in SnapDrive, you may get an error saying:

{% highlight text %}
A timeout of 120 secs elapsed while waiting for volume arrival notification from the operating system.
{% endhighlight %}

This generally occurs when you have created a LUN using the `NetApp System Manager` first, and then tried to connect to it in `SnapDrive`.

The reason for this error message is that when SnapDrive is connecting to a disk, it expects to see a formatted partition being connected – and when this doesn't occur, it doesn't know what to do and times out.

![NetApp LUN Timeout](/images/posts/netapp_lun_timeout/01_netapp_lun_timeout.png "NetApp LUN Timeout")

* TOC
{:toc}

## Clean Up The Mess

As the LUN mapping has failed half way through, you need to remove the igroups that have been added to the LUN, otherwise you will receive a `The specified LUN /vol/xxx is already mapped to at least one initiator` error.

* Connect to the filer using telnet
* Run the following command:

{% highlight text %}
lun unmap /vol/YourLUNName youriGroupName
{% endhighlight %}

You will get a message back saying the igroup has been unmapped from the LUN.

## Connect The LUN Again

Open SnapDrive and start connecting to your LUN again, as the process starts to occur you will see the following status messages appearing:

![Connect the LUN Again](/images/posts/netapp_lun_timeout/02_connect_the_lun_again.png "Connect the LUN Again")

* At this point, open up **Computer Management** and then in the left navigation pane select **Disk Management**
* You will notice a new, unformatted drive appear. If you don't see it keep refreshing until it appears

![LUN Attached as Un-formatted Volume](/images/posts/netapp_lun_timeout/03_new_unformated_volume.png "LUN Attached as Un-formatted Volume")

* You will need to be quick (you only have 120 seconds until time out remember), so initialize the disk

![LUN Attached as Un-formatted Volume](/images/posts/netapp_lun_timeout/04_initialize_the_disk.png "LUN Attached as Un-formatted Volume")

* After the disk has initialized, right click on the empty drive space and choose `Format` and follow the prompts to format the disk and specify the drive letter you choose when connecting to the LUN using SnapDrive

![Format the Disk](/images/posts/netapp_lun_timeout/05_format_the_disk.png "Format the Disk")

SnapDrive should now discover the disk and map it correctly. To prevent this occurring in the future, create your LUN's using the SnapDrive interface instead of NetApp System Manager, as SnapDrive will format the LUN on creation.
