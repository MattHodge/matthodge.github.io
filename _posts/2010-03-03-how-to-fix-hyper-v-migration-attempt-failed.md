---
layout: post
title:  "How To Fix Hyper-V Migration Attempt Failed"
date:   2010-03-03 13:37:00
comments: false
modified: 2016-02-05
---

If you are running Windows Server 2008 R2 with Hyper-V and have a SAN available, Hyper-V Live migration is a great way to make your Virtual Machines highly available. There is nothing like being able to migrate your VM’s to another node in the cluster and shutdown a host to perform maintenance – I still get excited every time I live migrate a machine!

In some cases though, your Live Migration attempt may fail. Unfortunately the Failover Cluster Manager nor the event viewer sheds much light on the reason in most cases.

![Failed Hyper-V Migration Attempt](/images/posts/hyperv_migration_attempt_failed/01_failed_migration_attempt.png "Failed Hyper-V Migration Attempt")

![Failed Migration Event Logs](/images/posts/hyperv_migration_attempt_failed/02_failed_event_logs.png "Failed Migration Event Logs")

To troubleshoot a problem with virtual no error message is difficult, and my process was to do a Google on the event log error and try my luck with the different forum posts from the TechNet Forums and other such sites. For this reason I decided to make a checklist to go through to try and work out the reason behind the failed live migration. Good luck!

* TOC
{:toc}

## The 'None Of My Hyper-V Servers Are Live Migrating' Checklist

* In the **Failover Cluster Manager** do your clustered networks have a status of `Up`?
* Do you have enough RAM free on the server you are trying to Live Migrate to? The best way to check this is System Center Virtual Machine Manager. If you do not have this, checking task manager on the destination server will show you how much RAM you have free.
* In the **Failover Cluster Manager** do your **Clustered Shared Volume(s)** has a status of `Online`?
* In the **Failover Cluster Manager** do your clustered **Storage** volumes have a status of `Online`?
* Are you able to ping the node you are trying to Live Migrate to and visa-versa?
* Did you create the Hyper-V network adapters **before** creating a cluster? If not you should destroy your cluster and start again reading the Microsoft guide [Using Hyper-V and Failover Clustering](https://technet.microsoft.com/en-us/library/cc732181%28WS.10%29.aspx)
* Are the names of your Hyper-V Virtual Adapters all the same?

## The 'Migration Attempt Failed but other VM's migrate fine' Checklist

These step is meant for people that are having a single Virtual Machine failing to Live Migrate for example; `VM1` and `VM2` are hosted on `NODE1`. `VM1` is successfully Live Migrated to `NODE2` but `VM2` says `Migration Attempt Failed`.

*	Make sure no CD/DVD or image file is mounted in the Virtual Machine you are trying to migrate.
* Reset the permission on the folder on the SAN that your Virtual Machine resides in to **Everyone** has full permission (this is not a security issue as the SAN network should be on a separate IP to your workstations and only Hyper-V administrators should have permission to access the cluster and its nodes)
* This is a weird one: Do you have any notes in the `Name` field under the **Management** section of the settings for the Virtual Machine? If so remove it and try Live Migration again. I personally have had this issue and removing the notes in here (that I never entered) worked.

![Weird Notes in the Notes Field](/images/posts/hyperv_migration_attempt_failed/03_notes_in_the_name_field.png "eird Notes in the Notes Field")

* If you have Anti-Virus software installed on one of your nodes, ensure it is set not to scan the `C:\ClusterStorage` folder

Hopefully using these checklists assists you in fixing your Hyper-V Live Migration problems. I will endeavour to update this list as I discover more reasons that Live Migration might fail.
