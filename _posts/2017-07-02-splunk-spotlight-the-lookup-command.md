---
layout: post
title: Splunk Spotlight - The Lookup Command
date: 2017-07-02T13:37:00.000Z
comments: true
description: We take a close look at the Splunk lookup command, including examples of where you might use it for enriching your logs and using CIDR matching.
---

Splunk is an amazing logging aggregation and searching tool. Even though I've been using it a few months now, I feel like I am just scratching the surface of what it can do.

My company recently switch over from the ELK stack (ElasticSearch, LogStash and Kibana) as we were moving to the cloud, with a focus on using managed services. The ELK stack is awesome, but it can be a pain to administer and extend. We were finding we spent more time administering our log collection pipeline as opposed to getting value from the logs it was storing.

I thought I would start a series of posts called "Splunk Spotlight" where I focus on a single feature or command inside Splunk and show some examples of how it can be used.

* TOC
{:toc}

# Getting Splunk Setup

The free edition of Splunk allows you to store 500mb/day. You can find a comparison of features [here](https://www.splunk.com/en_us/products/splunk-enterprise/free-vs-enterprise.html). You can use the free version for these examples.

The easiest way to play around with Splunk is to use Docker. I have setup a repository at [https://github.com/MattHodge/splunk](https://github.com/MattHodge/splunk) which I will keep updated with demo data files as I add more posts.

Make sure you have installed [docker-compose](https://docs.docker.com/compose/install/).

{% highlight bash %}
# Clone the repo
git clone git@github.com:MattHodge/splunk.git

# Enter the docker directory
cd splunk/docker

# Run docker compose to bring up the containers
docker-compose up -d
{% endhighlight %}

Once the container comes up, open up a browser and go to [http://localhost:8000/](http://localhost:8000/).

Enter the username of `admin` and password of `changeme` and you will be presented with the first Splunk screen.

# The Lookup Command

Time for our first Splunk command!

The Splunk `lookup` commands allows you to use data from an external source to enrich the data you already have in Splunk.

The *external source* can be one of the following:

* a CSV file (`CSV lookup`)
* a Python script or binary executable (`External lookup`)
* a Splunk key/value store collection (`KV Store lookup`)
* a Keyhole Markup Language (KML) file used for map lookups (`Geospatial lookup`)

You can find the full documentation for the `lookup` command on the Splunk documentation page [here](http://docs.splunk.com/Documentation/Splunk/6.6.2/Knowledge/Aboutlookupsandfieldactions).

# Getting Test Data

I have created some fake test data from [Mockaroo](https://www.mockaroo.com/) for the examples. All IP's and Data used in the examples is fake.

Let's upload it to our Splunk instance.

## Create an Index

To keep our Splunk installation clean, let's first create an index to store the demo data.

* Click on **Settings** > **Indexes** > **New Index**
* For index name put `webshop_demo` and choose the App of `Search & Reporting`. Leave everything else default
* Repeat the process and create another index called `iis_access_logs`

## Upload the Data

To add data do your Splunk installation:

* Click on **Settings** > **Add Data**

![Add Data to Splunk](/images/posts/splunk-lookup-command/splunk-add-data.png)

* Click on **Upload** and choose the `retail_orders.csv` from the `demo_data` folder that you cloned
* Keep clicking next until you get to **Input Settings**. On this page enter a **Host Field Value** of `webnode1` and for **index** choose `webshop_demo`
* Repeat again for `iis_access_logs.csv` from the `demo_data`. Put them into the `iis_access_logs` index

## Viewing the Data

* Once the data has imported, head back to the search page, and find your data by using the following query:

{% highlight bash %}
# View the webshop order logs
index="webshop_demo"

# View the webshop access logs
index="iis_access_logs"
{% endhighlight %}

![View the demo data](/images/posts/splunk-lookup-command/webshop_demo_data.png)

Now that we have loaded our data, let's look at some examples.

# Example 1 - Customer Ordering Data

In this example, let's pretend we have an online shop. We instrument the code of the shop to send a log message to Splunk every time someone makes a purchase.

Inside the product ordering code, we have access to the following values that we can log:

* `date_of_purchase` The date and time an order was created
* `order_id` The ID of the order that was created
* `payment_method_id` The ID of the payment method the customer used, which references a column in our retail applications database to a payment method table. Payment methods include Paypal, CreditCard, Gift Card, Cash on Delivery and Debit Card
* `order_amount` The value of the order in dollars

Our goal is to create a dashboard to see the types of payment methods the orders created with.

Unfortunately, in our code we don't have the text values for *Payment Method*.

We don't want to do a database query to find them out every time we send an order as this would slow down our ordering process. We also don't want to hard code the payment method names in our code, as new payment might be added at any time by our billing team.

![Payment Method IDs](/images/posts/splunk-lookup-command/payment_method_id.png)

This is where a Splunk `lookup` can help.

* First, let's create a lookup table. We will use a CSV file for this:

{% highlight csv %}
payment_method_id,payment_method_name
1,PayPal
2,Visa
3,Mastercard
4,Cash on Delivery
5,Gift Card
{% endhighlight %}

* Save this CSV file as `payment_methods.csv`.

We need to upload this file to Splunk so it can use it to do lookups on our data.

* Go to **Settings** > **Lookups**

![Splunk Lookup Settings](/images/posts/splunk-lookup-command/lookup_settings.png)

* Choose **Lookup table files** and click on **New**
* Upload the `payment_methods.csv` file and give it a destination file name of `payment_methods.csv`

![Add Splunk Lookup File](/images/posts/splunk-lookup-command/add_lookup_file.png)

> :white_check_mark: **Tip:** To allow other people to use the lookup file, you will need to edit the permission to make it shared in App.

Next we need to let Splunk know how to use the lookup file we added, and how it can use it to match and enrich fields in our searches.

* Go back to the **Lookups** screen and choose **Lookup definitions**
* Click **New**
* Enter the name of `payment_method`, choose the **Type** of `File-based` and choose the **Lookup file** of `payment_methods.csv`

Splunk will detect the supported fields in the CSV file.

![Splunk Supported Fields from CSV](/images/posts/splunk-lookup-command/supported_fields_from_csv.png)

* Go back to the search page and let's try out our lookup

To perform the lookup, the command looks like this:

{% highlight bash %}
index="webshop_demo" | lookup payment_methods.csv payment_method_id
{% endhighlight %}

![Splunk Lookup Search](/images/posts/splunk-lookup-command/splunk_lookup_search.png)

Splunk is matching `payment_method_id` from our lookup csv file and adding the additional field `payment_method_name`. This allows us to use the name of the payment method instead of the value when we make our dashboards.

![Pie Chart Visualization](/images/posts/splunk-lookup-command/piechart_visualization.png)

# Example 2 - Web Server Access Logs

In this example, let's pretend we have been asked by security to make a report of the top 5 IP Addresses that accessed the `login.html` page on our web application. We need to get this from our web server access logs.

Easy you say!

{% highlight bash %}
index="iis_access_logs" cs_uri_stem="/login.html" | top limit=5 c_ip
{% endhighlight %}

You run this query and give security the results.

![Top 5 Hits on Login with our IPs](/images/posts/splunk-lookup-command/top5_with_our_ips.png)

Security comes back and says "can you make this again, but this time not include any of our own IP addresses?". You look at the top 5 and realize that 3 of them are actually coming from the companies two office locations. This makes sense as many employees use the web application, but we need a way to filter those out.

The public IP ranges for those offices are:

* `200.13.37.0/24`
* `243.200.10.0/28`

This time, lets use a `KV Store lookup`. You can create and update a [KV store using the Splunk REST API](http://dev.splunk.com/view/webframework-developapps/SP-CAAAEZG), but we will use a Splunk Addon to manage the KV Store via the Web UI.

* Click the cog to manage Apps

![Manage Splunk Apps](/images/posts/splunk-lookup-command/manage_splunk_apps.png)

* Click on **Browse More Apps** and search for `Lookup File Editor` and click install (you will need a Splunk.com account to install apps)
* Once the App is installed, head back home and you will see **Lookup Editor** in the apps list on the left. Click on it

![Lookup Editor App](/images/posts/splunk-lookup-command/lookup_editor.png)

* Click **Create New Lookup** > **KV Store Lookup**
* Give the KV store the name of `office_ips` and store it in the `Search & Reporting` app
* Use two fields:
  * `c_ip` which will be a string, to match the IP address of the IIS logs
  * `isOfficeIP` which will be a string saying either `true` or `false` depending on if the IP address is from one of the offices. (You can use a `boolean` for this but you will get `1` or a `0`, and I prefer the `true`, `false`)

![Lookup Creation](/images/posts/splunk-lookup-command/lookup-edit.png)

* Enter the two CIDR IP ranges and enter `true` in the `isOfficeIP` column.

![Edit the lookup table](/images/posts/splunk-lookup-command/edit_lookup_table.png)

* Go to **Settings** > **Lookups** and choose **Lookup definitions**
* Click **New**
* Choose the `search` app as the destination. Enter the name of `office_ips`, choose the **Type** of `KV Store`. Enter the **Collection name** of `office_ips`. In **Supported Fields** enter `c_ip,isOfficeIP`.
* Expand out **Advanced Options** and set **Minimum Matches** to `1`. Set **Default Matches** to `false`. This is so that even if the IP doesn't match, we will get an `isOfficeIP` of `false`. For **Match Type** enter `CIDR(c_ip)`. This will make Splunk match the IPs in the web server logs `c_ip` field to the CIDR ranges we used in the KV Store.

![Configure Splunk KV Store Lookup](/images/posts/splunk-lookup-command/confgiure-kv-store-lookup-splunk.png)

* Go back to the search page and let's try out our lookup

{% highlight bash %}
index="iis_access_logs" cs_uri_stem="/login.html" | lookup office_ips c_ip
{% endhighlight %}

Once you do the search, you will see a new field is added to the events showing which IP's are in the office ranges, and which are not.

![IP's in Office CIDR Range](/images/posts/splunk-lookup-command/ips_in_office_cidr_range.png)

We can filter by only hits inside our office IP range.

{% highlight bash %}
index="iis_access_logs" | lookup office_ips c_ip | search isOfficeIP=true
{% endhighlight %}

![Show IP's in our CIDR Range](/images/posts/splunk-lookup-command/only_ips_in_our_cidr_range.png)

Now we can finally give the security team the report they want.

{% highlight bash %}
index="iis_access_logs" cs_uri_stem="/login.html" | lookup office_ips c_ip | search isOfficeIP=false | top limit=5 c_ip
{% endhighlight %}

![Top 5 IP's on Login Page but not from the offices](/images/posts/splunk-lookup-command/top5_ips_on_login_no_office_ips.png)

# Conclusion

The Splunk `lookup` command is a wonderful way to enrich your data after it has already been collected. It can help make your searches and dashboards more useful by giving you contextual information. You can also use the powerful CIDR matching functionality to group IP addresses and search based on things like offices or VLANs.

If you want more information, go and check out the [documentation](http://docs.splunk.com/Documentation/Splunk/latest/Knowledge/Addfieldsfromexternaldatasources) over on the Splunk Docs site.

Would you like to know when more of these "Splunk Spotlight?" posts come out? Make sure you follow me on Twitter [@MattHodge](https://twitter.com/MattHodge) and I will post new articles there.
