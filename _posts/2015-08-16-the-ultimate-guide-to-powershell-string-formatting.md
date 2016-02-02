---
layout: post
title:  "The Ultimate Guide to PowerShell String Formatting"
date:   2015-08-06 13:37:00
comments: true
modified: 2016-02-02
---

When scripting with PowerShell, you will come to a point where you need to work with strings that contain variables. Depending on your situation, there are several methods you can use when formatting strings in PowerShell. This blog will walk through these options.

We will start simple and ramp up the complexity.

{% gist b314defaf95de1e80494 %}

We run into our first problem here. We are using a dollar sign in our string, so PowerShell thinks it is a variable and tries to insert it into a string.

![Write-Output with Variable Doesn't Display](/images/posts/ps_string_formatting/write-output-with-variable-error.png "Write-Output with Variable Doesn't Display")

As we haven't set a variable for $6, the string we get back is incorrect. To fix this problem, we need to escape the dollar sign so PowerShell leaves it alone. In PowerShell, escapes are done using a backtick (\`).

{% gist 87d36cb92f98d8bec927 %}

Let's create a hash table for our beer and try and use its properties in a string.

Once everything is installed, we will be customizing and setting up the following:

{% gist 7a5b2b205f3895a7949b %}

We run into another problem where – PowerShell isn't extracting the value of the properties inside our hash table.

![Write-Output Hashtable Variable Not Showing](/images/posts/ps_string_formatting/write-output-hashtable-variable-not-showing.png "Write-Output Hashtable Variable Not Showing")

We can wrap variables with properties inside some special tags to force PowerShell to return our variable inside the string, for example `$($myItem.price)`. This is called a sub-expression:

{% gist 4eaf9fd4d38fc556956e %}

What if we need to use single quotes inside our string?

{% gist 61e9f5ddd0ce2b5cc082 %}

This works fine, but if we wanted to use double quotes, we have 2 options. We can escape using a backtick (\`) or escape using double-double quotes.

{% gist 2aa9598d70274058a1fe %}

With the above knowledge we can handle all string formatting situations. Here a final complex example of creating a HTML page out of a PowerShell string. It includes single and double quotes, dollar signs and hash tables and uses the format command.

{% gist 4a16275707261e81b36c %}

The beauty of PowerShell is that you can get the job done in multiple ways. You most likely would not mix all of them in a single script, but instead stick with the one you feel most comfortable with.

The above examples should make it a breeze to do even the most complicated string formatting. Happy PowerShell-ing!

Thanks to [@Jaykul](https://twitter.com/Jaykul) for reviewing this post.
