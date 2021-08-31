# errorLogging

This is a simple service that watches specified files and POSTs updates thereto to a Slack app (or any open webhook URL).

## Download the Linux bin here:
https://photographywesterncape.com/errorLogging

## Usage
Run the command to start a file-watcher process:

`./errorLogging --url|-u https://hooks.slack.com/services/.../... --files|-f /path/to/file/one.log@ERROR#ERROR,CRITICAL#ERROR /path/to/file/two.log@INFO#INFO,DEBUG#WARNING ...`

## Note
If filters are specified for a file (by using the "@" symbol after the file path followed by a comma-separated list of filters), then only file writes containing those particular strings will be reported. For example, "production.ERROR".

The "#" symbol after the filter string changes the formatting of messages associated with that filter. 

There are three formats available:
INFO (formats messages with a neutral border);
ERROR or CRITICAL (formats messages with a red border); and, 
DEBUG or WARNING (formats messsages with a yellow border).


