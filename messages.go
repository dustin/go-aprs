package aprs

import (
	"fmt"
	"regexp"
	"strings"
)

// A Message is a message from an APRS address to another.
type Message struct {
	Sender    Address
	Recipient Address
	Body      string
	ID        string
	Parsed    bool
}

// Message returns the message from an Frame frame.
func (a Frame) Message() Message {
	// Find source of third party
	for a.Body.Type().IsThirdParty() && len(a.Body) > 11 {
		a = ParseFrame(string(a.Body[1:]))
	}

	rv := Message{}
	if a.Body.Type().IsMessage() {
		if len(a.Body) < 12 {
			return rv
		}
		rv.Sender = a.Source
		rv.Recipient = AddressFromString(strings.TrimSpace(string(a.Body[1:10])))
		parts := strings.SplitN(string(a.Body[11:]), "{", 2)
		rv.Body = parts[0]
		if len(parts) > 1 {
			rv.ID = parts[1]
		}

		rv.Parsed = true
	}
	return rv
}

func (m Message) String() string {
	idstring := ""
	if m.ID != "" {
		idstring = "{" + m.ID
	}
	return fmt.Sprintf(":%-9s:%s%s", m.Recipient.String(),
		m.Body, idstring)
}

var (
	ackPattern = regexp.MustCompile(`^ack([A-z0-9]{1,5})`)
	blnPattern = regexp.MustCompile(`^:BLN[0-9]     :(.*)`)
	annPattern = regexp.MustCompile(`^:BLN[A-Z]     :(.*)`)
)

// IsACK returns true if this message is an acknowledgment to another message.
func (m Message) IsACK() bool {
	return ackPattern.MatchString(m.Body)
}

// IsBulletin returns true if the message represents a bulletin.
func (m Message) IsBulletin() bool {
	return blnPattern.MatchString(m.String())
}

// IsAnnouncement returns true if the message represents an announcement.
func (m Message) IsAnnouncement() bool {
	return annPattern.MatchString(m.String())
}
