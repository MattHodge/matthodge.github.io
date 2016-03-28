---
layout: post
title: "Windows Metric Dashboards with InfluxDB and Grafana"
date: 2016-03-28 13:37:00
comments: false
description: Using InfluxDB and Grafana to display Windows performance counters on beautiful dashboards.
---

.. to be completed ..

* TOC
{:toc}

## Requirements

You will need a Linux machine which will host the InfluxDB and Grafana installations. I will be using Ubuntu 14.04 x64 for this blog.

## Preparing the Ubuntu Machine

There is nothing special that needs to be performed on the Ubuntu server before installing InfluxDB or Grafana. Just make sure all the packages are up to date:

{% highlight bash %}
sudo apt-get update
sudo apt-get upgrade
{% endhighlight %}

The other thing I would recommend doing is setting the time zone of the Ubuntu server to UTC. It is a good idea to standardize on UTC as the time zone for all your metrics. InfluxDB uses UTC so stick to it. (You can read about some of the struggles when you don't use UTC [here](https://github.com/influxdata/influxdb/issues/2074)).

## Install InfluxDB

InfluxDB is a metric storage system.. put more info.

Download and install the InfluxDB .deb

{% highlight bash %}
cd /tmp
wget https://s3.amazonaws.com/influxdb/influxdb_0.11.0-1_amd64.deb
sudo dpkg -i influxdb_0.11.0-1_amd64.deb

# Start the service
sudo service influxdb start
{% endhighlight %}

InfluxDB listens on 2 main ports:
* TCP port `8083` is used for InfluxDB’s Admin panel
* TCP port `8086` is used for client-server communication over InfluxDB’s HTTP API

Go to `http://Your-Linux-Server-IP:8083` in the browser and confirm you can access the InfluxDB admin panel:

![InfluxDB Admin Panel](/images/posts/influxdb_grafana_windows/influxdbadminpanel.png "InfluxDB Admin Panel")

# Install Grafana

Grafana is a dashboard.. put more info.

Download and install the Grafana .deb

{% highlight bash %}
cd /tmp
wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.6.0_amd64.deb
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

# Install the Telegraf Client

Telegraf is an agent written in Go for collecting metrics from the system it's running on, or from other services, and writing them into InfluxDB or other outputs.

As the Windows agent is still in an experimental phase, head over to its GitHub page at https://github.com/influxdata/telegraf to grab the latest version.

At the time of writing the latest version could be found at http://get.influxdb.org/telegraf/telegraf-0.11.1-1_windows_amd64.zip.

Extract the zip file into a directory, I used `C:\telegraf`.

Inside you will see 2 files:
* `telegraf.exe` - this is the application. It is written in Go which compiles nicely into a single `.exe` file
* `telegraf.conf` - all the configuration options for telegraf



## Configure Telegraf

Open the `telegraf.conf` file in a text editor - I would recommend one which supports [TOML](https://github.com/toml-lang/toml) syntax highlighting such as [Atom](https://atom.io/).

The Windows version of telegraf has a configuration file setup to  collect some common Windows Performance Counters by default, so we do not need to change a lot.

The first thing we will change is the collection interval which I will bring down to 5 seconds. This lives under the `[agent]` section of the config:

{% highlight toml %}
[agent]
  interval = "5s"
{% endhighlight %}

P.S: I have removed the comments from the configuration file in these examples.

Next, under the `[[outputs.influxdb]]` section, we need to update the `urls` option to point to our InfluxDB server at `http://Your-Linux-Server-IP:8086`.

{% highlight toml %}
[[outputs.influxdb]]
  urls = ["http://Your-Linux-Server-IP:8086"]
{% endhighlight %}

Save the `telegraf.conf`.

To run telegraf, open and then we will start telegraf with the following command:

{% highlight powershell %}
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

The default database for the telegraf agent is `telegraf`. The Grafana form will not let you save unless you enter a User and Password, so just enter in something random as we have not configured credentials just yet.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_datasource.png "Configure Grafana Data Source")

To display our data, we will need to create a dashboard. Select `Home` from the top menu and click `New`.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_create_dashboard.png "Configure Grafana Data Source")

In the new dashboard page you will see a little green rectangle over on the left, click it and choose `Add Panel` > `Graph`.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_add_panel.png "Configure Grafana Data Source")

Click on the `Metrics` tab, and down the bottom drop down the data source dropdown choose the data source we added, `InfluxDB`.

In the data selection section, choose From `win_cpu` and match the rest of the fields up to the image below to get a graph of the CPU usage.

![Configure Grafana Data Source](/images/posts/influxdb_grafana_windows/grafana_cpu_graph.png "Configure Grafana Data Source")
