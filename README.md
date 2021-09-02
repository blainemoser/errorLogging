# errorLogging

This is a service that watches specified files and POSTs relevant writes thereto to a Slack app's webhook.

## Download the Linux bin here:
https://photographywesterncape.com/errorLogging

## Usage
Run the command to start a file-watcher process:

`./errorLogging --url|-u https://hooks.slack.com/services/.../... --files|-f /path/to/file/one.log@ERROR#ERROR,CRITICAL#ERROR /path/to/file/two.log@INFO#INFO,DEBUG#WARNING --suppress|-s "error that you want to suppress" "something else that you want to suppress" ...`

## Filters
If filters are specified for a file (by using the "@" symbol after the file path followed by a comma-separated list of filters), then only file writes containing those particular strings will be posted. 

For example, by specifying the file path with the filter "production.ERROR,production.CRITICAL" ("/path/to/file@production.ERROR,production.CRITICAL") the programme will only post writes that contain either "production.ERROR" and/or "production.CRITICAL".

## Formatting
Use the "#" symbol after the filter string to format associated Slack posts. 

For example specifying "/path/to/file@production.WARNING#WARNING,production.CRITICAL#ERROR" applies the WARNING format to Slack messages for writes containing the string "production.WARNING", and the ERROR format for those containing "production.CRITICAL" (see below).

There are three formats available:

# INFO 
(formats messages with a neutral border);

# ERROR or CRITICAL 
(formats messages with a red border); and,

# DEBUG or WARNING 
(formats messages with a yellow border).

## Suppress Messages
Use the suppress (--suppress|-s) option to ignore certain messages. If any file write contains one or more of the specified strings, it will not be posted.

