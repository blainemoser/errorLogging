# errorLogging

This is a service that watches specified files and POSTs relevant writes thereto to a Slack app's webhook.

## Download the Linux bin here:
https://photographywesterncape.com/errorLogging

## Usage
Run the command to start a file-watcher process:

`./errorLogging --url|-u https://hooks.slack.com/services/.../... --files|-f /path/to/file/one.log@ERROR#ERROR,CRITICAL#ERROR /path/to/file/two.log@INFO#INFO,DEBUG#WARNING --suppress|-s "error that you want to suppress" "something else that you want to suppress" ...`

## Filters
If filters are specified for a file (by using the "@" symbol after the file path followed by a comma-separated list of filters), then only file writes containing those particular strings will be posted. For example, by specifying the file path with the filter "production.ERROR,production.CRITICAL" ("/path/to/file@production.ERROR,production.CRITICAL") the programme will only post writes that contain one of the strings "production.ERROR" and "production.CRITICAL".

The "#" symbol after the filter string changes the formatting of messages associated with that filter. For example specifying "/path/to/file@production.ERROR#WARNING,production.CRITICAL#ERROR" means posts of writes containing the string "production.ERROR" will have the WARNING formatting, whereas those containing the string "production.CRITICAL" will have ERROR formatting (see below).

There are three formats available:

INFO (formats messages with a neutral border);

ERROR or CRITICAL (formats messages with a red border); and,

DEBUG or WARNING (formats messages with a yellow border).

## Suppress certain messages
Use the suppress (--suppress|-s) option to specify text that you want suppressed. If any file write contains any of the strings specified, it will not be posted.

