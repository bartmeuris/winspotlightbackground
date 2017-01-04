# Windows spotlight background manager

This program copies the images Windows 10 fetches from Bing to use as lock-screen when "Windows Spotlight" is selected as source in the control panel. This allows you to use these images as normal background. 

The best way to use this is to create a scheduled task which is triggered upon login and unlock.

The program can:

* copy only landscape or portrait images. By default it only copies portrait images.
* ignores small images (< 150x150px)
* only copy files with a specific width or height
* remove duplicates from the target directory, based on image size, filesize and file hash. It keeps the most recent file.
* completely manage the target directory (remove all files not matching the filter criteria)

# Commandline options:

* `-target`: Specify a target folder. By default it is the `%USERPROFILE%\Pictures\Spotlight` folder. If this directory does not exists, it is created.
* `-landscape`: Copy landscape images. This is the default. To disable, specify `-landscape=false`.
* `-portrait`: Copy portrait images. Useful when you have vertical displays. Disabled by default.
* `-width <width>`: Only copy images with a widht matching the specified number.
* `-height <height>`: Only copy images with a height matching the specified number.
* `-copysmall`: Copy small (<150x150) images anyway. Default: false.
* `-targetdedup`: Remove duplicate files from the target folder (see `-target`).
* `-targetvalidateremove`: Removes all images in the target directory not matching the specified filters. **This is dangerous and can lead to data loss.**
* `-logfile <filename>`: Log messages to the specified file. By default, nothing is logged.

# Recommended usage

## Create a scheduled task

* Open `Task Scheduler` in Windows 10.
* In the `Task Scheduler`, optionally create a new folder there to keep things clean, and select the folder.
* Right click and choose `Create New Task`.
  * Give it a name
  * Run as the user you want it to provide the backgrounds for
  * Select to only run when the user is logged in
  * Go to the `Triggers` tab and add 2 triggers:
    * `At logon`, select `Specific user` and pick the correct user account
    * `On workstation unlock`, and again, select `Specific user` and pick the correct user account.
  * Go to the `Actions` tab, and create a new action:
    * browse for the `winspotlightbackground.exe` program
    * Fill in the additional arguments (see Commandline options)
  * Verify the other settings of the scheduled task to meet your demands

It is recommended to test this task at least once by right-clicking it and selecting `Run`, so the target folder is created if it didn't exist, and some initial images are being copied. Verify this after running the task by opening an explorer and navigating to the `%USERPROFILE%\Pictures\Spotlight` folder (if you didn't override the target folder with the `-target` command line parameter).

## Configure the desktop background

This is the easy part, to ensure you actually use the newly copied images:

* right click your desktop background
* select `Personalize`
* in the `Background` section, select `Slideshow` from the dropdown menu
* select `Browse`, enter `%USERPROFILE%\Pictures\Spotlight` and press enter. Then click `Choose this folder`.

Close the settings window and you should be good to go!

