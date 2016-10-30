package aprs

import (
	"testing"
)

const (
	MESSAGE  = "KG6HWF-9>APDR12,TCPIP*,qAC,T2SPAIN2::KG6HWF   :testing notifications{10"
	MESSAGE2 = "K7FED-10>APJI23,WR6ABD*:}KG6HWE>APOA00,TCPIP,K7FED-10*::KG6HWF   :yo{AB}07"
	ACKED    = "KG6HWF-5>APDR13,TCPIP*,qAC,T2PERTH::KG6HWF   :ack01}1"
)

func TestMessage(t *testing.T) {
	v := ParseFrame(MESSAGE)
	msg := v.Message()

	if !msg.Parsed {
		t.Fatalf("Couldn't parse %v as a message", v)
	}
	if msg.Sender.String() != "KG6HWF-9" {
		t.Fatalf("Didn't find the sender: %v", msg.Recipient)
	}
	if msg.Recipient.String() != "KG6HWF" {
		t.Fatalf("Didn't find the receipient: %v", msg.Recipient)
	}
	if msg.Body != "testing notifications" {
		t.Fatalf("Didn't get the message: %#v from %#v", msg.Body, v.Body)
	}
	if msg.ID != "10" {
		t.Fatalf("Expected msg id 10, got %v", msg.ID)
	}
}

func TestAcked(t *testing.T) {
	v := ParseFrame(ACKED)
	msg := v.Message()

	if !msg.IsACK() {
		t.Fatalf("Expected %v to be interpreted as an ACK", ACKED)
	}
}

func TestBrokenMessage(t *testing.T) {
	a := Frame{Body: ":"}
	msg := a.Message()
	if msg.Parsed {
		t.Fatalf("Expected to fail to parse broken message: %v", msg)
	}
}

func TestThirdParty(t *testing.T) {
	v := ParseFrame(MESSAGE2)
	if !v.Body.Type().IsThirdParty() {
		t.Fatalf("This should be third party traffic: %#v", v.Body)
	}
	msg := v.Message()

	if !msg.Parsed {
		t.Fatalf("Couldn't parse %v as a message", v)
	}
	if msg.Sender.String() != "KG6HWE" {
		t.Fatalf("Incorrect sender: %v", v.Source)
	}
	if msg.Recipient.String() != "KG6HWF" {
		t.Fatalf("Didn't find the receipient: %v", msg.Recipient)
	}
	if msg.Body != "yo" {
		t.Fatalf("Didn't get the message: %#v from %#v", msg.Body, v.Body)
	}
	if msg.ID != "AB}07" {
		t.Fatalf("Expected msg id AB}07, got %v", msg.ID)
	}
}

func TestMessageEncoding(t *testing.T) {
	exp := ":KG6HWF   :yo{AB}07"
	m := Message{Sender: AddressFromString("KG6HWE"),
		Recipient: AddressFromString("KG6HWF"),
		Body:      "yo",
		ID:        "AB}07",
	}
	if m.String() != exp {
		t.Fatalf("Expected %v, got %v", exp, m.String())
	}
}
