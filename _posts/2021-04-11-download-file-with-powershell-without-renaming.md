---
layout: post
title: Downloading a file with PowerShell without specifying its name
date: 2021-04-11-T00:00:00.000Z
comments: false
description: How to download a file using PowerShell using the name of the download instead of specifying it yourself.
---

I'm going to show you how to download a file from the web using PowerShell _without_ having to specify the filename on download.

This definately falls into the category of "why is this so hard in PowerShell?". Looks like some other folks agree: [github.com/PowerShell/issues/11671](https://github.com/PowerShell/PowerShell/issues/11671).

Usually you would need to use `-Outfile` to download a file, for example:

```powershell
Invoke-WebRequest -Uri "https://go.microsoft.com/fwlink/?linkid=2109047&Channel=Stable&language=en&consent=1" -Outfile "MicrosoftEdgeSetup.exe"
```

How does the browser know what file to give the download though?

![Edge Download](images/posts/download-file-powershell/edge-setup-download.png)

Let's find out:

```powershell
$download = Invoke-WebRequest -Uri "https://go.microsoft.com/fwlink/?linkid=2109047&Channel=Stable&language=en&consent=1"

$download.Headers
```

With that we can see the name of the download:

![Download](images/posts/download-file-powershell/download_name.png)

To get the file name out:

```powershell
$content = [System.Net.Mime.ContentDisposition]::new($download.Headers["Content-Disposition"])
$content
```

![Show Filename](images/posts/download-file-powershell/show_filename.png)

Now we _could_ take that file name and run another `Invoke-WebRequest` with a `-Outfile` parameter, but that would involve downloading the entire file again.

Let's save the contents of our `$download` variable to disk.

```powershell
$fileName = $content.FileName

$file = [System.IO.FileStream]::new($fileName, [System.IO.FileMode]::Create)
$file.Write($download.Content, 0, $download.RawContentLength)
$file.Close()
```

Now we have `MicrosoftEdgeSetup.exe` saved.

![Finally](images/posts/download-file-powershell/finally.gif)

Here is the `Save-Download` function which makes this process easier:

{% gist c102c75d65420852fe8424ba8e75ba25 %}