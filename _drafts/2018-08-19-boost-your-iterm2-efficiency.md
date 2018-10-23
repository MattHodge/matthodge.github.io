---
layout: post
title: Boost your iTerm2 Efficiency
date: 2018-08-19T13:37:00.000Z
comments: false
description: How to improve efficiency when working with iTerm2 on Mac OS.
---

* TOC
{:toc}

## Add a hotkey for quick access

* Go to **iTerm2 | Preferences | Keys**
* Under **Hotkey**, enter your prefered hotkey to quickly show the iTerm2 window.

## Jump forwards and backwards between words

Have a lot of text in a command and holding down the arrow key taking to long? Let's enable jumping forwards and backwards between words.

* Go to **iTerm2 | Preferences | Profiles**
* Choose your profile and select the **Keys** header
* Select the radio button for **Left ⌥ Key: Esc+**

![Set iTerm2 Escape Key](/images/posts/iterm2/set-profile-escape-key.png)

* Click the **+** sign to add a new keyboard shortcut
* Click on **Click to Set** and enter **⌥ ←** and for the **Action** choose **Send escape sequence**
* In the **Esc+** box enter **b** an click **OK**

![Set the escape key left](/images/posts/iterm2/set-escape-key-left.png)

* Do the same for **⌥ →**, instead entering **f** in the **Esc+** box

![Set the escape key right](/images/posts/iterm2/set-escape-key-right.png)

Now you can quickly skip over words.

![Skip over words](/images/posts/iterm2/skip-over-words.gif)

## Keyboard shortcuts

Here are some handy keyboard shortcuts:

Details | Command
--- | ---
Clear the screen (just like typing `clear`) | **⌘ r**
Split the screen | **⌘ d**
Jump between the split screens | **⌘ [** and **⌘ ]**
