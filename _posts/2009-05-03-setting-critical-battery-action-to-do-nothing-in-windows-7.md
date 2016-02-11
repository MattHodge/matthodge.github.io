---
layout: post
title:  "Setting the Critical Battery Action to 'Do Nothing' In Windows 7"
date:   2009-05-03 13:37:00
comments: false
modified: 2016-02-10
---
For some reason, Windows 7 does not allow setting the Critical Battery Action to `Do Nothing`.

You may wish to use this if the Advanced Configuration and Power Interface (ACPI) in your laptop is reporting your battery charge level incorrectly.

My poor old IBM ThinkPad T30 has this issue, so I want to force Windows not to power off the laptop when the battery is low.

The tool we will be using that comes with Windows 7 and Vista is `powercfg.exe`.

* Activate the power scheme you want to modify.
* Open an Administrative command prompt.
* Enter: `powercfg -setdcvalueindex SCHEME_CURRENT SUB_BATTERY BATACTIONCRIT 0`

Your current power scheme will show `Battery | Critical Battery Action | On Battery : Do Nothing`, even though the option is not available in the drop down box.

Now you can let your battery run dry properly!
