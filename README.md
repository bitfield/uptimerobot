[![Go Reference](https://pkg.go.dev/badge/github.com/bitfield/uptimerobot.svg)](https://pkg.go.dev/github.com/bitfield/uptimerobot)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitfield/uptimerobot)](https://goreportcard.com/report/github.com/bitfield/uptimerobot)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)
![Tests](https://github.com/bitfield/uptimerobot/actions/workflows/test.yml/badge.svg)

# uptimerobot

`uptimerobot` is a Go library and command-line client for the [Uptime Robot](https://uptimerobot.com/) website monitoring service. It allows you to search for existing monitors, delete monitors, create new monitors, and also inspect your account details and alert contacts.

## Installing the command-line client

To install the client binary, run:

```
go get -u github.com/bitfield/uptimerobot
```

## Running the command-line client in Docker

To use the client in a Docker container, run:

```
docker container run bitfield/uptimerobot
```

## Using the command-line client

To see help on using the client, run:

```
uptimerobot -h
```

## Setting your API key

To use the client with your Uptime Robot account, you will need the Main API Key for the account. Go to the [Uptime Robot Settings page](https://uptimerobot.com/dashboard.php#mySettings) and click 'Show/hide it' under the 'Main API Key' section.

There are three ways to pass your API key to the client: in a config file, in an environment variable, or on the command line.

### In a config file

The `uptimerobot` client will read a config file named `.uptimerobot.yaml` (or `.uptimerobot.json`, or any other extension that Viper supports) in your home directory, or in the current directory.

For example, you can put your API key in the file `$HOME/.uptimerobot.yaml`, and `uptimerobot` will find and read it automatically (replace `XXX` with your own API key):

```yaml
apiKey: XXX
```

### In an environment variable

`uptimerobot` will look for the API key in an environment variable named UPTIMEROBOT_API_KEY:

```
export UPTIMEROBOT_API_KEY=XXX
uptimerobot ...
```

(For historical reasons, the variable can also be named `UPTIMEROBOT_APIKEY`.)

### On the command line

You can also pass your API key to the `uptimerobot` client using the `--apiKey` flag like this:

```
uptimerobot --apiKey XXX ...
```

## Testing your configuration

To test that your API key is correct and `uptimerobot` is reading it properly, run:

```
uptimerobot account
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
uptimerobot contacts
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
uptimerobot search www.example.com
ID: 780689017
Name: Example.com website
URL: https://www.example.com/
Status: Up
Type: HTTP
```

(Use `uptimerobot monitors` to list all existing monitors.)

If there are no monitors found matching your search, the exit status of the command will be 1. Otherwise it will be 0. (If you're checking whether a monitor already exists before creating it, try the `ensure` command instead.)

## Deleting monitors

Note the ID number of the monitor you want to delete, and run `uptimerobot delete`:

```
uptimerobot delete 780689017
Monitor ID 780689017 deleted
```

## Pausing or starting monitors

Note the ID number of the monitor you want to pause, and run `uptimerobot pause`:

```
uptimerobot pause 780689017
Monitor ID 780689017 paused
```

To resume a paused monitor, run `uptimerobot start` with the monitor ID:

```
uptimerobot start 780689017
Monitor ID 780689017 started
```

## Creating a new monitor

Run `uptimerobot new URL NAME` to create a new monitor:

```
uptimerobot new https://www.example.com/ "Example.com website"
New monitor created with ID 780689018
```

To create a new monitor with alert contacts configured, use the `-c` flag followed by a comma-separated list of contact IDs, with no spaces:

```
uptimerobot new -c 0102759,2053888 https://www.example.com/ "Example.com website"
New monitor created with ID 780689019
```

## Ensuring a monitor exists

Sometimes you want to create a new monitor only if a monitor doesn't already exist for the same URL. This is especially useful in automation.

To do this, run `uptimerobot ensure URL NAME`:

```
uptimerobot ensure https://www.example.com/ "Example.com website"
Monitor ID 780689018 ensured
```

If the monitor doesn't already exist, it will be created.

You can use the `-c` flag to add alert contacts, just as for the `uptimerobot new` command.

## Checking the version number

To see what version of the command-line client you're using, run `uptimerobot version`.

## Viewing debug output

When things aren't going quite as they should, you can add the `--debug` flag to your command line to see a dump of the HTTP request and response from the server. This is helpful if you want to report problems with the client, for example.

## Using the Go library

If the command-line client doesn't do quite what you need, or if you want to use Uptime Robot API access in your own programs, import the library using:

```go
import "github.com/bitfield/uptimerobot/pkg"
```

Create a new `Client` object by calling `uptimerobot.New()` with an API key:

```go
client = uptimerobot.New(apiKey)
```

Once you have a client, you can use it to call various Uptime Robot API features:

```go
monitors, err := client.AllMonitors()
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
	ID           int64  `json:"id,omitempty"`
        FriendlyName string `json:"friendly_name"`
        URL          string `json:"url"`
        ...
}
```

For example, to delete a monitor, find the ID of the monitor you want to delete, and pass it to `DeleteMonitor()`:

```go
if err := client.DeleteMonitor(780689017); err != nil {
        log.Fatal(err)
}
```

To call an Uptime Robot API verb not implemented by the `uptimerobot` library, you can use the `MakeAPICall()` method directly, passing it some suitable JSON data:

```go
r := uptimerobot.Response{}
data := []byte(fmt.Sprintf("{\"id\": \"%d\"}", m.ID))
if err := client.MakeAPICall("deleteMonitor", &r, data); err != nil {
    log.Fatal(err)
}
fmt.Println(r.Monitor.ID)
```

The API response is returned in the `Response` struct. If the call fails, `MakeAPICall()` will return the error message. Otherwise, the requested data will be available in the appropriate field of the `Response` struct:

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

For example, when creating a new monitor, the ID of the created monitor will be returned as `r.Monitor.ID`.

If things aren't working as you expect, you can use the debug facility to dump the raw request and response data from every API call. To do this, set the environment variable `UPTIMEROBOT_DEBUG`, which will dump debug information to the standard output, or set `client.Debug` to any `io.Writer` to send output to that writer.

Here's an example of the debug output shown when creating a new monitor:

```http
POST /v2/newMonitor HTTP/1.1
Host: api.uptimerobot.com
User-Agent: Go-http-client/1.1
Content-Length: 221
Content-Type: application/json
Accept-Encoding: gzip

{
  "alert_contacts": "0335551_0_0-2416450_0_0",
  "api_key": "XXX",
  "format": "json",
  "friendly_name": "Example check",
  "port": 443,
  "type": 1,
  "url": "https://www.example.com"
}

HTTP/2.0 200 OK
Access-Control-Allow-Origin: *
Cf-Ray: 505422654b04dbf3-LHR
Content-Type: application/json; charset=utf-8
Date: Mon, 12 Aug 2019 17:22:57 GMT
Etag: W/"33-NlNt8dOhQvno31TtQYsI0xTJ9w"
Expect-Ct: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
Server: cloudflare
Set-Cookie: __cfduid=d9ec99b8a777d9f806956432718fb5c81565630577; expires=Tue, 11-Aug-20 17:22:57 GMT; path=/; domain=.uptimerobot.com; HttpOnly
Vary: Accept-Encoding

{"stat":"ok","monitor":{"id":783263671,"status":1}}
```

## Bugs and feature requests

If you find a bug in the `uptimerobot` client or library, please [open an issue](https://github.com/bitfield/uptimerobot/issues). Similarly, if you'd like a feature added or improved, let me know via an issue.

Not all the functionality of the Uptime Robot API is implemented yet.

Pull requests welcome!
