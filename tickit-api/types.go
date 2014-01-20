package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GenericAPIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type TicketsAPIResponse struct {
	Since   time.Time         `json:"since"`
	Last    time.Time         `json:"last"`
	Message string            `json:"message"`
	Success bool              `json:"success"`
	Tickets map[string]Ticket `json:"tickets"`
}

type ScansAPIResponse struct {
	Since   time.Time       `json:"since"`
	Last    time.Time       `json:"last"`
	Message string          `json:"message"`
	Success bool            `json:"success"`
	Scans   map[string]Scan `json:"scans"`
}

type ManifestsAPIResponse struct {
	Last      time.Time  `json:"last"`
	Message   string     `json:"message"`
	Success   bool       `json:"success"`
	Manifests []Manifest `json:"manifests"`
}

type RecordAPIResponse struct {
	Message      string   `json:"message"`
	Success      bool     `json:"success"`
	UnsavedUUIDs []string `json:"unsaved_uuids"`
}

type Ticket struct {
	TicketNumber string `json:"ticket_number"`
	TicketHolder string `json:"ticket_holder"`
	EventName    string `json:"event_name"`
	ItemName     string `json:"item_name"`
}

type Scan struct {
	UUID         string    `json:"uuid"`
	TicketNumber string    `json:"ticket_number"`
	Time         time.Time `json:"scan_time"`
	Location     string    `json:"scan_location"`
	Count        int       `json:"scan_count"`
	Reversal     bool      `json:"scan_reversal"`
	OrderID      int       `json:"-"`
	ItemID       int       `json:"-"`
	OK           bool      `json:"scan_valid"`
}

// A list of active event items (i.e. ticket types) requested before loading individual tickets by item_id
type Manifest struct {
	EventID    int    `json:"event_id"`
	EventTitle string `json:"event_title"`
	ItemID     int    `json:"item_id"`
	ItemTitle  string `json:"item_title"`
}

func GenerateTicketNumber(orderID int, itemID int, sequence int, letter string) (ticketNumber string, err error) {
	ticketNumber = fmt.Sprintf("%d-%d-%d-%s", orderID, itemID, sequence, letter)
	return
}

func ParseTicketNumber(tn string) (orderID int, itemID int, sequence int, letter string, err error) {
	parts := strings.SplitN(tn, "-", 4)

	if len(parts) < 4 {
		err = errors.New("Invalid ticket number")
		return
	}
	orderID, _ = strconv.Atoi(parts[0])
	itemID, _ = strconv.Atoi(parts[1])
	sequence, _ = strconv.Atoi(parts[2])
	letter = parts[3]

	return
}
