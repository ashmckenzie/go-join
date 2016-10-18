package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "strings"

  "github.com/urfave/cli"
)

// APIURL ..
const APIURL string = "https://joinjoaomgcd.appspot.com/_ah/api/messaging/v1/sendPush"

// DEBUG ...
var DEBUG = false

// VERBOSE ...
var VERBOSE = false

func send(c *cli.Context) {
  hc := http.Client{}

  APIKey := c.GlobalString("api-key")
  deviceID := c.GlobalString("device-id")
  URL := c.GlobalString("url")
  iconURL := c.GlobalString("icon")
  title := c.GlobalString("title")

  if len(title) == 0 {
    title, _ = os.Hostname()
  }

  body := c.Args().Get(0)
  if body == "-" {
    bytes, _ := ioutil.ReadAll(os.Stdin)
    body = strings.Trim(string(bytes), "\n")
  }

  form := url.Values{}
  form.Add("apikey", APIKey)
  form.Add("deviceId", deviceID)
  form.Add("title", title)
  form.Add("text", body)
  form.Add("icon", iconURL)
  form.Add("url", URL)

  url := fmt.Sprintf("%s?%s", APIURL, form.Encode())

  req, err := http.NewRequest("POST", url, nil)
  if err != nil {
    log.Fatal(err)
  }

  _, err = hc.Do(req)
  if err != nil {
    log.Fatal(err)
  }
}

func validateParams(c *cli.Context) error {
  if len(c.GlobalString("api-key")) == 0 {
    return cli.NewExitError("ERROR: API key is empty", 1)
  }

  if len(c.GlobalString("device-id")) == 0 {
    return cli.NewExitError("ERROR: Device ID empty", 2)
  }

  if len(c.Args().Get(0)) == 0 {
    return cli.NewExitError("ERROR: Body is empty", 3)
  }

  return nil
}

func main() {
  app := cli.NewApp()

  app.Name = "join"
  app.Usage = "Send a push using Join"
  app.Version = os.Getenv("VERSION")

  app.Flags = []cli.Flag{
    cli.BoolFlag{
      Name:        "verbose",
      Usage:       "Verbose mode",
      EnvVar:      "VERBOSE",
      Destination: &VERBOSE,
    },
    cli.BoolFlag{
      Name:        "debug",
      Usage:       "Debug mode",
      EnvVar:      "DEBUG",
      Destination: &DEBUG,
    },
    cli.StringFlag{
      Name:   "api-key, a",
      Usage:  "API key to use",
      EnvVar: "JOIN_API_KEY",
    },
    cli.StringFlag{
      Name:   "device-id, d",
      Usage:  "Device ID to send message to",
      EnvVar: "JOIN_DEVICE_ID",
    },
    cli.StringFlag{
      Name:   "title, t",
      Usage:  "Title of the message",
      EnvVar: "JOIN_TITLE",
    },
    cli.StringFlag{
      Name:   "url, u",
      Usage:  "URL to use",
      EnvVar: "JOIN_URL",
    },
    cli.StringFlag{
      Name:   "icon, i",
      Usage:  "ICON URL to use",
      EnvVar: "JOIN_ICON",
    },
  }

  app.Action = func(c *cli.Context) error {
    err := validateParams(c)
    if err != nil {
      log.Fatal(err)
    }
    send(c)
    return nil
  }

  app.Run(os.Args)
}
