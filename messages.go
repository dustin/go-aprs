package aprs

import (
	"strings"
)

type Message struct {
	Sender    Address
	Recipient Address
	Body      string
	Id        string
	Parsed    bool
}

func (a APRSMessage) Message() (rv Message) {
	// Find source of third party
	for a.Body.Type().IsThirdParty() && len(a.Body) > 11 {
		a = ParseAPRSMessage(string(a.Body[1:]))
	}

	if a.Body.Type().IsMessage() {
		rv.Sender = a.Source
		rv.Recipient = AddressFromString(strings.TrimSpace(string(a.Body[1:10])))
		parts := strings.SplitN(string(a.Body[11:]), "{", 2)
		rv.Body = parts[0]
		if len(parts) > 1 {
			rv.Id = parts[1]
		}

		rv.Parsed = true
	}
	return
}
