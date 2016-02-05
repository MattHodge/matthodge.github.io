---
layout: post
title:  "Resending Reporting Services Email Subscriptions When Troubleshooting"
date:   2010-08-12 13:37:00
comments: true
modified: 2016-02-05
---

Troubleshooting Microsoft Reporting Services email subscriptions can be a pain. How many times have you created a new subscription a few minutes in advance and not received it, then you are unsure if you set the schedule time correctly, or if you pressed the save button, or if the fix you made to reporting services didn't work? It's a time consuming process, but thankfully, there is a better way.

![Don't Set Reports in the Future!](/images/posts/resend_reporting_services_email/01_reports_in_the_future.png "Don't Set Reports in the Future!")

Instead, you can run some SQL commands on your reporting services to trigger the running and resending of this report.

* Use **SQL Server Management Studio** and connect to the **database engine** of your reports server
* Click the **New Query** button
* Run this SQL query to list of all the reports with schedules:

{% highlight sql %}
SELECT
sj.[name] AS [Job Name],
c.[Name] AS [Report Name],
c.[Path],
su.Description,
su.EventType,
su.LastStatus,
su.LastRunTime
FROM msdb..sysjobs AS sj INNER JOIN ReportServer..ReportSchedule AS rs
ON sj.[name] = CAST(rs.ScheduleID AS NVARCHAR(128)) INNER JOIN
ReportServer..Subscriptions AS su
ON rs.SubscriptionID = su.SubscriptionID INNER JOIN
ReportServer..[Catalog] c
ON su.Report_OID = c.ItemID
{% endhighlight %}

* From the **Results** pane, determine the job name of the report you want to trigger. Right click on the the job Guid in the `JobName` column and click copy

![List the Reports](/images/posts/resend_reporting_services_email/02_list_the_reports.png "List the Reports")

* Click the **New Query** button again to open a blank query window
* Run this SQL query, replacing `YourJobNameHere` with your `Job Name` retrieved from the last step

{% highlight sql %}
USE [msdb]
EXEC sp_start_job @job_name = 'YourJobNameHere'
{% endhighlight %}

* When you execute the query, the Message window should say `Job 'GUID'` started successfully:

![Start Job Manually](/images/posts/resend_reporting_services_email/03_start_job_manually.png "Start Job Manually")

If you don't receive the report – then you know you didn't fix the initial problem, but now at least, you have a fast way to resend the report each time you change a Reporting Services / SMTP setting!
