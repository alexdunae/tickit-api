package main

import (
	"testing"
)

func Test_ParseTicketNumber(t *testing.T) {
	_, _, _, _, err := ParseTicketNumber("")
	if err == nil {
		t.Error("Blank ticket number should yield error")
	}

	_, _, _, _, err = ParseTicketNumber("123-45-1")
	if err == nil {
		t.Error("Incomplete ticket number should yield error")
	}

	orderID, itemID, sequence, letter, err := ParseTicketNumber("123-45-9-CC")
	if orderID != 123 {
		t.Errorf("Order ID incorrect, was %v", orderID)
	}

	if itemID != 45 {
		t.Errorf("Item ID incorrect, was %v", itemID)
	}

	if sequence != 9 {
		t.Errorf("Sequence incorrect, was %v", sequence)
	}

	if letter != "CC" {
		t.Errorf("Letter incorrect, was %v", letter)
	}

	if err != nil {
		t.Error("Should not yield error")
	}
}

func Test_GenerateTicketNumber(t *testing.T) {
	const expected = "91012-1234-52-WW"
	ticketNumber, err := GenerateTicketNumber(91012, 1234, 52, "WW")
	if err != nil {
		t.Error("Should not yield error")
	}

	if ticketNumber != expected {
		t.Errorf("Incorrect ticket number, was %v, expected %v", ticketNumber, expected)

	}
}
