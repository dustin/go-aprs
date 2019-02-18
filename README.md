go-aprs - APRS library for Golang
======


# Installation

```shell
go get github.com/dustin/go-aprs
```


# Usage



## Read-Only stream

For those interested in Read-Only stream, the `NewROAPRS()` method allows an easy way
to receive APRS messages.


```go

package main

import (
    "fmt"
    "log"
    "github.com/dustin/go-aprs/aprsis"
)

func main() {
    client, err := aprsis.NewROAPRS()

    if err != nil {
        log.Fatal("login", err)
    }

    defer client.Close()

    go client.ManageConnection()

    for frame := range client.GetIncomingMessages() {
        fmt.Println(frame.Source.Call, frame.Dest.Call, frame.Body.Type().String())
    }
}
```

## Read-Write stream

Please replace `CALLSIGN-4` with your callsign and `12345` with your passcode.
`r/10.5/10.5/500` is APRS filter expression: latitude, longitude, radius (km).


```go

package main

import (
    "fmt"
    "log"
    "github.com/dustin/go-aprs/aprsis"
)

func main() {
    client, err := aprsis.NewAPRS("CALLSIGN-4", "12345", "r/10.5/10.5/500")

    if err != nil {
        log.Fatal("login", err)
    }

    defer client.Close()

    packet := "CALLSIGN-4>APZU5,TCPIP*,qAC,T2PARIS:=1030.28N/01030.20E-"

    err = client.SendRawPacket(packet)
    if err != nil {
        log.Fatalf("Send %+v\n", err)
    }

    go client.ManageConnection()

    for frame := range client.GetIncomingMessages() {
        fmt.Println(frame.Source.Call, frame.Dest.Call, frame.Body.Type().String())
    }
}
```
