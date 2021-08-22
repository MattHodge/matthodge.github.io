---
layout: post
title: Securing Home Assistant with Cloudflare
date: 2020-07-27-T00:00:00.000Z
comments: true
description: A guide to securing your Home Assistant installation with the Cloudflare CDN.
---

*Updated: Aug 22nd, 2021 due to a [HTTP Proxy breaking change](https://github.com/home-assistant/core/pull/51839) in Home Assistant.*

I've just started using [Home Assistant](https://www.home-assistant.io/) through building my own smart garage door opener that I could control using my phone.

It's an amazing piece of open source software, and very easy to get setup locally, but I wanted to expose it to the internet so I could see the status of my garage door when away from the house using the [Home Assistant App](https://www.home-assistant.io/integrations/mobile_app/).

I am running Home Assistant Core with Docker on my home server, and was a little concerned about opening my home server up to the internet, especially one where you could open a door into my house remotely. Home Assistant has had a [very good history](https://www.cvedetails.com/vulnerability-list/vendor_id-17232/product_id-41425/Home-assistant-Home-assistant.html) when it comes to security vulnerabilities in their software, but I wanted to be as careful as I could.

This article I will describe using [Cloudflare's free plan](https://www.cloudflare.com/plans/) to protect remote access to Home Assistant.

This article assumes:

* You have your own domain name
* You have your domain setup to use Cloudflare nameservers

If you want to register a domain, I recommend [Namecheap](https://www.namecheap.com/). You can then set it up in Cloudflare using [these docs](https://www.namecheap.com/support/knowledgebase/article.aspx/9607/2210/how-to-set-up-dns-records-for-your-domain-in-cloudflare-account). If you already have a domain, you can follow [the docs here](https://support.cloudflare.com/hc/en-us/articles/201720164-Creating-a-Cloudflare-account-and-adding-a-website), to set it up in Cloudflare.

> 📢 Want to know when more posts like this come out? [Follow me on Twitter: @MattHodge](https://twitter.com/matthodge) 📢

* TOC
{:toc}

## Choose a new port for Home Assistant

The default port for Home Assistant (`8123`) is not supported when proxied through Cloudflare.

At the time of writing, the [supported ports for HTTPS](https://support.cloudflare.com/hc/en-us/articles/200169156-Identifying-network-ports-compatible-with-Cloudflare-s-proxy) are as follows:

* 443
* 2053
* 2083
* 2087
* 2096
* 8443

Choose a port from the list, and configure the Home Assistant [HTTP integration](https://www.home-assistant.io/integrations/http/) in the `configuration.yaml`:

```yaml
# Example configuration.yaml entry
http:
  server_port: 8443 # Use your chosen port
```

Restart Home Assistant and confirm you can still access it locally. Update the port forward on your router so you can access your Home Assistant instance over the internet.

## Setup a subdomain for your Home Assistant

In Cloudflare, create a subdomain in the **DNS** tab for your domain.

* Click `+ Add Record`
* Choose type `A` and add your subdomain (I used `hass` in my example below)
* For the `IPv4` address, enter the IP address of your home internet connection
* Ensure that the `Proxy Status` option is set to `Proxied`. This is how we can leverage Cloudflare to protect our Home Assistant instance

![Create Cloudflare Subdomain](images/posts/securing-home-assistant-with-cloudflare/subdomain.png)

If you don't have a static IP address on your home internet connection, you can use the [Home Assistant Cloudflare](https://www.home-assistant.io/integrations/cloudflare/) addon to keep it up to date.

You should now be able to access your Home Assistant using the subdomain via Cloudflare.

## Setup an SSL Certificate

Cloudflare provides free SSL certificates automatically. Try hitting `https://<subdomain>.<domain>:<port>` and you should be accessing Home Assistant over SSL.

This provides an encrypted connection from your web browser to Cloudflare, but the connection from Cloudflare to your server is still un-encrypted.

![Cloudflare to Home Assistant not encrypted](images/posts/securing-home-assistant-with-cloudflare/backend-not-encrpyted.png)

To encrypt communication between Cloudflare and Home Assistant, we will use an Origin Certificate.

In Cloudflare, got to the **SSL/TLS** tab:

* Click `Origin Server`
* Click `Create Certificate`

![Cloudflare to Home Assistant not encrypted](images/posts/securing-home-assistant-with-cloudflare/create-origin-cert.png)

* Enter the subdomain that the Origin Certificate will be generated for

![Origin Certificate Generation Dialog](images/posts/securing-home-assistant-with-cloudflare/cert-generation-dialog.png)

In the next dialog you will be presented with the contents of two certificates.

* Copy the contents of `Origin Certificate` and save it in a file called `origin.pem`
* Copy the contents of `Private Key` and save it in a file called `privkey.pem`

Update your `configuration.yaml` with the following, replacing the path with something accessible by your Home Assistant installation:

```yaml
# Example configuration.yaml entry
http:
  server_port: 8443 # From the previous step
  ssl_certificate: /certificate/path/origin.pem
  ssl_key: /certificate/path/privkey.pem
```

Restart Home Assistant and access it with `https://<subdomain>.<domain>:<port>`, which should be the same as before, but will now be encrypted end to end.

![End to End Encrypted Access to Home Assistant](images/posts/securing-home-assistant-with-cloudflare/end-to-end-encrypted.png)

You can also optionally enable `Full (strict)` encryption.

![Enable full encryption](images/posts/securing-home-assistant-with-cloudflare/enable-full-encryption.png)

## Blocking Traffic Not Originating From Cloudflare

We now have our encrypted traffic going through Cloudflare, but if someone gets our home IP address, they can go around Cloudflare and hit our Home Assistant directly.

![Attacker avoiding Cloudflare](images/posts/securing-home-assistant-with-cloudflare/attacker-avoiding-cloudflare.png)

To prevent this, you can configure your firewall to only allow traffic to Home Assistant to Cloudflare IP addresses. Cloudflare [lists all their IP addresses here](https://www.cloudflare.com/ips/).

I am using **ufw** on Ubuntu, and used Ansible to configure the firewall on the home server running Home Assistant, but you can do this manually in whatever firewall you are using.

```yaml
{% raw %}
- name: home assistant local 8443/tcp
    ufw:
    rule: allow
    src: "192.168.1.0/24" # your LAN range
    port: "8443"
    direction: in
    proto: tcp

# Example Ansible configuration to allow only Cloudflare IPs into Home Assistant
- name: home assistant remote from cloudflare ips (ipv4) 8443
    ufw:
    rule: allow
    src: '{{ item }}'
    port: "8443"
    direction: in
    proto: tcp
    loop: "{{ lookup('url', 'https://www.cloudflare.com/ips-v4', wantlist=True, headers={'User-Agent':'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11'}) }}" # Without a header this request is blocked.
{% endraw %}
```

Now only Cloudflare IPs will be able to access your Home Assistant

## Allow Proxy Requests to Home Assistant

Home Assistant provides some built in protection for proxy servers (for example CloudFlare) access to your Home Assistant installation as of version [2021.7](https://www.home-assistant.io/blog/2021/07/07/release-20217/).

To allow CloudFlare to work as a proxy, modify your `http` config (part of your `configuration.yaml`):

```yaml
# Add use_x_forwarded_for
use_x_forwarded_for: true
# Add the Cloudflare IPs as trusted proxies https://www.cloudflare.com/ips-v4
trusted_proxies:
  - 173.245.48.0/20
  - 103.21.244.0/22
  - 103.22.200.0/22
  - 103.31.4.0/22
  - 141.101.64.0/18
  - 108.162.192.0/18
  - 190.93.240.0/20
  - 188.114.96.0/20
  - 197.234.240.0/22
  - 198.41.128.0/17
  - 162.158.0.0/15
  - 104.16.0.0/13
  - 104.24.0.0/14
  - 172.64.0.0/13
  - 131.0.72.0/22
```

## Using the Cloudflare Firewall

Even though we now have Cloudflare protecting our Home Assistant, anyone on the internet can still access it and try logging in:

![Home Assistant Login Attempt](images/posts/securing-home-assistant-with-cloudflare/home-assistant-login-attempt.png)

To prevent this, we can the Cloudflare firewall to further restrict access.

In Cloudflare, got to the **Firewall** tab:

* Click `Firewall Rules`
* Click `Create a Firewall Rule`

![Create Firewall Rule](images/posts/securing-home-assistant-with-cloudflare/create-firewall-rule.png)

For example, I am only allowing connections to my Home Assistant from the Netherlands where I live:

![Create Firewall Rule](images/posts/securing-home-assistant-with-cloudflare/block-all-but-netherlands.png)

Keep in mind you may need to create some exceptions if you have incoming webhooks or other automation hitting your Home Assistant instance from the internet.

You can use the `Firewall Events` view in the Cloudflare console to troubleshoot this.

## Setup Two Factor Authentication

We have some good protections for our Home Assistant in place now, but it is a good idea to also enable one of the [Two Factor Authentication](https://www.home-assistant.io/docs/authentication/multi-factor-auth/) options Home Assistant provides.

## Conclusion

Following this guide, you will now have a fairly secure Home Assistant setup running on your home network. Happy automating!
