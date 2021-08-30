# errorLogging

This is a simple service that watches specified files and POSTs updates thereto to a Slack app (or any open webhook URL).

*Usage*
Run the command to start a file-watcher process:
`errorLogging --url|-u https://hooks.slack.com/services/.../... --files|-f /path/to/file/one.log@ERROR,CRITICAL /path/to/file/two.log@INFO ...`

If filters are specified for a file (by using the "@" symbol after the file path and entering a comma-separated list of filters) then only file writes containing those filters will be reported. 

