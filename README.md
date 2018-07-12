# uptimerobot

`uptimerobot` is a Go library and command-line client for the [UptimeRobot](https://uptimerobot.com/) website monitoring service. It allows you to search for existing monitors, delete monitors, create new monitors, and also inspect your account details and alert contacts.

## Installing the command-line client

To install the client binary, run:

```
go get -u github.com/bitfield/uptimerobot
```

## Using the command-line client

To see help on using the client, run:

```
uptimerobot -h
```

To use the client with your UptimeRobot account, you will need the Main API Key for the account. Go to the [UptimeRobot Settings page](https://uptimerobot.com/dashboard.php#mySettings) and click 'Show/hide it' under the 'Main API Key' section.

Copy the key to the clipboard and pass it to the `uptimerobot` client using the `--apiKey` flag like this (replace `XXX` with your own API key):

```
uptimerobot --apiKey XXX account
```

You should see your account details listed:

```
Email: j.random@example.com
Monitor limit: 300
Monitor interval: 1
Up monitors: 208
Down monitors: 2
Paused monitors: 0
```

If you get an error message, double-check you have the correct API key:

```
2018/07/12 16:04:26 API error: {
 "message": "api_key not found.",
 "parameter_name": "api_key",
 "passed_value": "XXX",
 "type": "invalid_parameter"
}
```

## Listing contacts

The `uptimerobot contacts` command will list your configured alert contacts by ID number:

```
uptimerobot --apiKey XXX contacts
ID: 0102759
Name: Jay Random
Type: 2
Status: 2
Value: j.random@example.com

ID: 2053888
Name: Slack
Type: 11
Status: 2
Value: https://hooks.slack.com/services/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
```

This will be useful when you create a new monitor, because you can add the contact IDs which should be alerted when the check fails (see 'Creating a new monitor' below).

## Listing or searching for monitors

Use `uptimerobot search` to list all monitors whose 'friendly name' or check URL match a certain string:

```
uptimerobot --apiKey XXX search www.example.com
ID: 780689017
Name: Example.com website
URL: https://www.example.com/
Type: HTTP
Subtype:
Keyword type: 0
Keyword value:
```

(Use `uptimerobot monitors` to list all existing monitors.)

If there are no monitors found matching your search, the exit status of the command will be 1. Otherwise it will be 0. This can be useful in automation scripts for checking whether or not a given monitor already exists.

## Deleting monitors

Note the ID number of the monitor you want to delete, and run `uptimerobot delete`:

```
uptimerobot --apiKey XXX delete 780689017
Monitor ID 780689017 deleted
```

## Creating a new monitor

Run `uptimerobot new URL NAME` to create a new monitor:

```
uptimerobot --apiKey XXX new https://www.example.com/ "Example.com website"
New monitor created with ID 780689018
```

To create a new monitor with alert contacts configured, use the `-c` flag followed by a comma-separated list of contact IDs, with no spaces:

```
uptimerobot --apiKey XXX new -c 0102759,2053888 https://www.example.com/ "Example.com website"
New monitor created with ID 780689019
```

## Using the Go library

If the command-line client doesn't do quite what you need, or if you want to use UptimeRobot API access in your own programs, import the library using:

```go
import "github.com/bitfield/uptimerobot/pkg"
```

Create a new `Client` object by calling `uptimerobot.New()` with an API key:

```go
client = uptimerobot.New(apiKey)
```

Once you have a client, you can use it to call various UptimeRobot API features:

```go
monitors, err := client.GetMonitors()
if err != nil {
        log.Fatal(err)
}
for _, m := range monitors {
        fmt.Println(m)
        fmt.Println()
}
```

Most API operations use the `Monitor` struct, which looks like this:

```go
type Monitor struct {
        ID           int64  `json:"id"`
        FriendlyName string `json:"friendly_name"`
        URL          string `json:"url"`
        ...
}
```

For example, to delete a monitor, first create a new `Monitor` variable and set its `ID` field to the ID of the monitor you want to delete. Then pass it to `DeleteMonitor()`:

```go
m := uptimerobot.Monitor{
        ID: 780689017,
}
new, err := client.DeleteMonitor(m)
if err != nil {
        log.Fatal(err)
}
fmt.Println(new.ID)
```

To call an UptimeRobot API verb not implemented by the `uptimerobot` library, you can use the `MakeAPICall()` method directly:

```go
r := uptimerobot.Response{}
p := uptimerobot.Params{
        "id": 780689017,
}
if err := client.MakeAPICall("deleteMonitor", &r, p); err != nil {
    log.Fatal(err)
}
fmt.Println(r.Monitor.ID)
```

The API response is returned in the `Response` struct. If the call fails, `MakeAPICall()` will return the  error message. Otherwise, the requested data will be available in the appropriate field of the `Response` struct:

```go
type Response struct {
        Stat          string         `json:"stat"`
        Account       Account        `json:"account"`
        Monitors      []Monitor      `json:"monitors"`
        Monitor       Monitor        `json:"monitor"`
        AlertContacts []AlertContact `json:"alert_contacts"`
        Error         Error          `json:"error"`
}
```

For example, when deleting a monitor, as in the above example, the ID of the deleted monitor will be returned as `r.Monitor.ID`.

## Bugs and feature requests

If you find a bug in the `uptimerobot` client or library, please [open an issue](https://github.com/bitfield/uptimerobot/issues). Similarly, if you'd like a feature added or improved, let me know via an issue.

Not all the functionality of the UptimeRobot API is implemented yet.

Pull requests welcome!
