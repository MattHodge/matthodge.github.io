---
layout: post
title: "Windows Metric Dashboards with InfluxDB and Grafana"
date: 2016-03-28 13:37:00
comments: false
description: Using InfluxDB and Grafana to display Windows performance counters on beautiful dashboards.
---

Understanding performance of your infrastructure is extremely important, especially when running production systems. There is nothing worse than a customer calling and saying they are experiencing slowness with one of their applications and you having no idea where to start looking.

In the [2014 State of DevOps survey](https://puppetlabs.com/sites/default/files/2014-state-of-devops-report.pdf) survey, one of the questions asked was **how is your organization notified of failure?**.

Here was the multiple choice question asked:

![how is your organization notified of failure](/images/posts/influxdb_grafana_windows/notification_of_failure.png "how is your organization notified of failure")

Through the survey, one of the top practices that correlated with performant IT teams was **Monitor system and application health**:

>Logging and monitoring systems make it easy to detect failures and
identify the events that contributed to them. Proactive monitoring of
system health based on threshold and rate-of-change warnings enables
us to preemptively detect and mitigate problems.

If you want some more information about performant DevOps teams and the methods they used to test teams, I recommend the talk [What We Learned from Three Years of Sciencing the Crap Out of DevOps](https://www.youtube.com/embed/cJVUtbSmXaM).

<iframe width="560" height="315" src="https://www.youtube.com/embed/cJVUtbSmXaM" frameborder="0" allowfullscreen></iframe>


Monitoring performance counters on Windows in any centralized manager way has always been tricky. In 2014 I wrote a PowerShell Module to send performance counters to Graphite which turned out to be pretty popular:

[Graphite-PowerShell-Functions](https://hodgkins.io/using-powershell-to-send-metrics-graphite) <a class="github-button" href="https://github.com/MattHodge/Graphite-PowerShell-Functions" data-icon="octicon-star" data-style="mega" data-count-href="/MattHodge/Graphite-PowerShell-Functions/stargazers" data-count-api="/repos/MattHodge/Graphite-PowerShell-Functions#stargazers_count" data-count-aria-label="# stargazers on GitHub" aria-label="Star MattHodge/Graphite-PowerShell-Functions on GitHub">Star</a>

Thankfully, things are getting easier. Let's take a look at using [InfluxDB](https://influxdata.com/) to store our metrics, [Telegraf](https://influxdata.com/time-series-platform/telegraf/) to tramsit the metrics and [Grafana](http://grafana.org/) do display them.


* TOC
{:toc}

## Requirements

You will need a Linux machine which will host the InfluxDB and Grafana installations. I will be using Ubuntu 14.04 x64 for this.

### Preparing the Ubuntu Machine

There is nothing special that needs to be performed on the Ubuntu server before installing InfluxDB or Grafana. Just make sure all the packages are up to date:

{% highlight shell %}
sudo apt-get update
sudo apt-get upgrade
{% endhighlight %}

### UTC Time

The other thing I would recommend is setting the time zone of the Ubuntu server to UTC. It is a good idea to standardize on UTC as the time zone for all your metrics. InfluxDB uses UTC so stick to it.

(You can read about some of the struggles when you don't use UTC [here](https://github.com/influxdata/influxdb/issues/2074)).

## Install InfluxDB

InfluxDB is an open source distributed time series database with no external dependencies. It's useful for recording metrics, events, and performing analytics. I recommend having a read of the [key concepts of InfluxDB](https://docs.influxdata.com/influxdb/v0.11/concepts/key_concepts/) over at their documentation page.

Let's download and install the InfluxDB `.deb`

{% highlight shell %}
cd /tmp
wget https://s3.amazonaws.com/influxdb/influxdb_0.11.0-1_amd64.deb
sudo dpkg -i influxdb_0.11.0-1_amd64.deb

# Start the service
sudo service influxdb start
{% endhighlight %}

InfluxDB listens on 2 main ports:

* TCP port `8083` is used for InfluxDB’s Admin panel
* TCP port `8086` is used for client-server communication over InfluxDB’s HTTP API

Once installed, go to `http://Your-Linux-Server-IP:8083` in the browser and confirm you can access the InfluxDB admin panel:

![InfluxDB Admin Panel](/images/posts/influxdb_grafana_windows/influxdbadminpanel.png "InfluxDB Admin Panel")

# Install Grafana

Grafana is a beautiful open source, metrics dashboard and graph editor. It can read data from multiple sources, for example Graphite, Elasticsearch, OpenTSDB, as well as InfluxDB. Take a look at the [Grafana live demo](http://play.grafana.org/) site to see what it can do.

First we will download and install the Grafana `.deb`. You can find the latest version over at [http://grafana.org/download/](http://grafana.org/download/).

{% highlight shell %}
cd /tmp
wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.6.0_amd64.deb

# Install required packages for Grafana
sudo apt-get install -y adduser libfontconfig
sudo dpkg -i grafana_2.6.0_amd64.deb

# Start the service
sudo service grafana-server start

# Configure Grafana to start at boot time
sudo update-rc.d grafana-server defaults 95 10
{% endhighlight %}

Grafana's web interface listens on TCP port `3000` by default.

Go to `http://Your-Linux-Server-IP:3000` in the browser and confirm you can access the InfluxDB admin panel:

![Grafana Login Panel](/images/posts/influxdb_grafana_windows/grafanalogin.png "Grafana Login Panel")

# Telegraf
Telegraf is an agent written in Go for collecting metrics from the system it's running on, or from other services, and writing them into InfluxDB or other outputs.

We will be using the `win_perf_counters` plugin for telegraf to collect Windows performance counters and send them over to InfluxDB. More information on the plugin can be found at the [telegraf GitHub page](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/win_perf_counters).

## Install the Telegraf Client

As the Windows agent is still in an experimental phase, head over to its GitHub page at [https://github.com/influxdata/telegraf](https://github.com/influxdata/telegraf) to grab the latest version.

At the time of writing the latest version could be found at [http://get.influxdb.org/telegraf/telegraf-0.11.1-1_windows_amd64.zip](http://get.influxdb.org/telegraf/telegraf-0.11.1-1_windows_amd64.zip).

Extract the zip file into a directory, I used `C:\telegraf`.

Inside you will see 2 files:

* `telegraf.exe` - this is the application. It is written in Go which compiles nicely into a single `.exe` file
* `telegraf.conf` - all the configuration options for telegraf

## Configure Telegraf

Open the `telegraf.conf` file in a text editor - I would recommend one which supports [TOML](https://github.com/toml-lang/toml) syntax highlighting such as [Atom](https://atom.io/).

The Windows version of telegraf has a configuration file setup to collect some common Windows performance counters by default, so we do not need to change very much for it to work.

The first thing we will change is the collection interval. This is how often the performance counters will be read. I am setting mine to 5 seconds. This configuration option is under the `[agent]` section:

{% highlight toml %}
[agent]
  interval = "5s"
{% endhighlight %}

Next, under the `[[outputs.influxdb]]` section, we need to update the `urls` option to point to our InfluxDB server at `http://Your-Linux-Server-IP:8086`.

{% highlight toml %}
[[outputs.influxdb]]
  urls = ["http://Your-Linux-Server-IP:8086"]
{% endhighlight %}

Save the `telegraf.conf` file.

To run telegraf, open and then we will start telegraf with the following command:

{% highlight yaml %}
C:\telegraf\telegraf.exe -config C:\telegraf\telegraf.conf
{% endhighlight %}

If all went well you should see telegraf starting to collect your metrics and send them over to InfluxDB.

![Starting Telegraf](/images/posts/influxdb_grafana_windows/startingtelegraf.png "Starting Telegraf")

## Troubleshooting

If you get an error saying `2016/03/28 19:48:01 toml: line 1: parse error` this is because you used standard old notepad and its line-endings broke things. Use a real text editor!

# Viewing the Data in Grafana

Now you have some metrics being sent into InfluxDB, you can use Grafana to view them.

Open up `http://Your-Linux-Server-IP:3000` and login using the default credentials:

* Username: `admin`
* Password: `admin`

## Configure a Data Source

Grafana needs to have a data source added so it knows where to look for the metrics.

Click on `Data Sources` on the left and then `Add new` at the top.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_configdatasource.png "Configure Grafana Data Source")

Choose the type `InfluxDB 0.9.x` for the data source and enter the URL for InfluxDB. Keep in mind that Grafana is running on the same box as InfluxDB, so you can just use `http://localhost:8086`.

Keep access as `proxy`.

The default database for the telegraf agent is `telegraf`. The Grafana form will not let you save unless you enter a User and Password, so just enter in something random as we have not configured any InfluxDB credentials.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_datasource.png "Configure Grafana Data Source")

To display our data, we will need to create a dashboard. Select `Home` from the top menu and click `New`.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_create_dashboard.png "Configure Grafana Data Source")

In the new dashboard page you will see a little green rectangle over on the left, click it and choose `Add Panel` > `Graph`.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_add_panel.png "Configure Grafana Data Source")

Click on the `Metrics` tab, and down on the bottom right of the page is the data source dropdown. Choose the data source we added, called `InfluxDB`.

In the data selection section, choose From `win_cpu` and match the rest of the fields up to the image below to get a graph of the CPU usage.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_cpu_graph.png "Configure Grafana Data Source")

You can read more about querying data from InfluxDB in Grafana in the [Grafana docs](http://docs.grafana.org/datasources/influxdb/).

# Wrapping Up

InfluxDB and telegraf provide an excellent and simple way to ship Windows performance counters off the server, and Grafana lets us display these metrics in beautiful dashboards.

Hopefully this starts you on your journey to graphing performance data for your systems.

Keep an eye out for another post shortly which will discuss some more advanced usage including using annotations on the graphs so you can correlate events in your infrastructure to system performance

<!-- Place this tag right after the last button or just before your close body tag. -->
<script async defer id="github-bjs" src="https://buttons.github.io/buttons.js"></script>
