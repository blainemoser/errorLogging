# errorLogging

This is a service that watches specified files and POSTs relevant writes thereto to a Slack app's webhook.

## Download the Linux bin here:
https://photographywesterncape.com/errorLogging

## Usage
Run the following command to start a file-watcher process:

#### ./errorLogging --url https://hooks.slack.com/services/.../... \
#### --files /path/to/file/one.log@ERROR#ERROR,CRITICAL#ERROR /path/to/file/two.log@INFO#INFO,DEBUG#WARNING \
#### --suppress "ignore this message" "also ignore this message"`

Abbreviations for the arguments are:
**-u|--url** (the Slack Webhook URL)
**-f|--files** (the file paths of the files to watch)
**-s|--suppress** (any messages that should be ignored)

## Filters
Using the "@" symbol after file paths limits writes to only those which contain those strings (use a comma-separated list to specify multiple filters). 

For example, by using these filters: 
#### production.ERROR,production.CRITICAL ("/path/to/file@production.ERROR,production.CRITICAL") 
...the programme will only post writes that contain either **production.ERROR** and/or **production.CRITICAL**.

## Formatting
Use the "#" symbol (after filter strings) to apply formatting to Slack messages (for writes that contain the filter string). 

For example specifying 
#### /path/to/file@production.WARNING#WARNING,production.CRITICAL#ERROR
...applies the WARNING format to Slack messages for writes containing the string **production.WARNING**, and the ERROR format for those containing **production.CRITICAL** (see below).

There are three formats available:

#### INFO 
(formats messages with a neutral border);

#### ERROR or CRITICAL 
(formats messages with a red border); and,

#### DEBUG or WARNING 
(formats messages with a yellow border).

#### Suppress Messages
Use the suppress (--suppress|-s) option to ignore certain messages. If any file write contains one or more of the specified strings, it will not be posted.

