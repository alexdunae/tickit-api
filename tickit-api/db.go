package main

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"strconv"
	"strings"
	"time"
)

func LoadManifests() (manifests []Manifest, err error) {
	manifests = []Manifest{}

	db, err := sql.Open("mymysql", config.DSN())

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(`
    SELECT events.id AS event_id, events.title AS event_title, items.id AS item_id, items.title AS item_title
			FROM items
			INNER JOIN events ON events.id = items.event_id
		WHERE
			events.enabled = 1
		AND
			# ensure there are some tickets
			(
				SELECT COUNT(*) FROM line_items
					INNER JOIN orders ON line_items.order_id = orders.id
				WHERE line_items.item_id = items.id AND line_items.affects_inventory = 1 AND orders.completed = 1
			) > 300 # TODO: should just be > 0
		ORDER BY events.position ASC, events.title ASC, items.position ASC, items.title ASC
    `)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		panic(err.Error())
	}

	var (
		eventID    int
		eventTitle string
		itemID     int
		itemTitle  string
	)

	for rows.Next() {
		err = rows.Scan(&eventID, &eventTitle, &itemID, &itemTitle)

		if err != nil {
			panic(err.Error())
		}

		manifest := Manifest{
			EventID:    eventID,
			EventTitle: eventTitle,
			ItemID:     itemID,
			ItemTitle:  itemTitle}

		manifests = append(manifests, manifest)
	}

	err = rows.Err()

	if err != nil {
		panic(err.Error())
	}

	return
}

func StoreExists(storeID int, key string) (exists bool) {
	db, err := sql.Open("mymysql", config.DSN())

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT id FROM stores WHERE id = ? AND checkin_key = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(storeID, key)

	if err != nil {
		panic(err.Error())
	}

	stores := make([]int, 0)
	var (
		foundID int
	)

	for rows.Next() {
		err = rows.Scan(&foundID)

		if err != nil {
			panic(err.Error())
		}

		stores = append(stores, foundID)
	}
	err = rows.Err()

	if err != nil {
		panic(err.Error())
	}

	log.Printf("verify %+v %d", stores, len(stores))

	return len(stores) == 1
}

func LoadTickets(itemIDs []int, since time.Time) (tickets map[string]Ticket, last time.Time) {
	tickets = map[string]Ticket{}
	last = time.Now().UTC()

	db, err := sql.Open("mymysql", config.DSN())

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(`
    SELECT orders.id, orders.name, line_items.item_id, items.title AS item_name,
    			 events.title AS event_name, line_items.quantity
      FROM line_items
      INNER JOIN orders ON orders.id = line_items.order_id
      INNER JOIN items ON items.id = line_items.item_id
      INNER JOIN events ON events.id = items.event_id
      WHERE
        affects_inventory = true AND
        item_id IN(?) AND
        completed = 1 AND
        orders.updated_at >= ?
    `)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	itemIDStrings := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		itemIDStrings[i] = strconv.Itoa(itemID)
	}

	rows, err := stmt.Query(strings.Join(itemIDStrings, ","), since)

	if err != nil {
		panic(err.Error())
	}

	var (
		orderID   int
		orderName string
		itemID    int
		itemName  string
		eventName string
		quantity  int
	)

	for rows.Next() {
		err = rows.Scan(&orderID, &orderName, &itemID, &itemName, &eventName, &quantity)

		if err != nil {
			panic(err.Error())
		}

		createTicketsFromRow(tickets, orderID, orderName, itemID, itemName, eventName, quantity)
	}
	err = rows.Err()

	if err != nil {
		panic(err.Error())
	}

	return
}

func createTicketsFromRow(tickets map[string]Ticket, orderID int, orderName string, itemID int, itemName string, eventName string, quantity int) {
	for sequence := 1; sequence <= quantity; sequence++ {
		ticketNumber, err := GenerateTicketNumber(orderID, itemID, sequence, "A")
		if err != nil {
			panic(err.Error())
		}
		tickets[ticketNumber] = Ticket{
			TicketNumber: ticketNumber,
			TicketHolder: orderName,
			EventName:    eventName,
			ItemName:     itemName}
	}
}

func LoadScans(itemIDs []int, since time.Time) (scans map[string]Scan, last time.Time) {
	scans = map[string]Scan{}
	last = time.Now().UTC()
	db, err := sql.Open("mymysql", config.DSN())

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// TODO: requires index on `item_id`, `scanned_at`
	stmt, err := db.Prepare(`
    SELECT * FROM (
    	SELECT ticket_number, scanned_at AS ScanTime, location AS ScanLocation, validated, updated_at
    	FROM check_ins
      WHERE
      	item_id IN(?)
    	ORDER BY scanned_at DESC
		) as c
		WHERE updated_at >= ?
		GROUP BY ticket_number`)
	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	itemIDStrings := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		itemIDStrings[i] = strconv.Itoa(itemID)
	}

	rows, err := stmt.Query(strings.Join(itemIDStrings, ","), since)

	if err != nil {
		panic(err.Error())
	}

	var (
		ticketNumber string
		scanTime     time.Time
		scanLocation string
		validated    bool
		updatedAt    time.Time
	)

	for rows.Next() {
		err = rows.Scan(&ticketNumber, &scanTime, &scanLocation, &validated, &updatedAt)

		scans[ticketNumber] = Scan{
			TicketNumber: ticketNumber,
			Time:         scanTime.UTC(),
			Location:     scanLocation,
			Count:        1,
			OK:           true,
			Reversal:     false}

		if err != nil {
			panic(err.Error())
		}
	}

	err = rows.Err()

	if err != nil {
		panic(err.Error())
	}
	return
}

func SaveScans(scans []Scan) (savedUUIDs []string, err error) {
	savedUUIDs = make([]string, 0)

	db, err := sql.Open("mymysql", config.DSN())

	if err != nil {
		return
	}
	defer db.Close()

	timestamp := time.Now().UTC()

	for _, scan := range scans {
		orderID, itemID, _, _, err := ParseTicketNumber(scan.TicketNumber)

		// invalid ticket number
		if err != nil {
			log.Printf("not saving invalid ticket number: %s", scan.TicketNumber)
			continue
		}

		scan.OrderID = orderID
		scan.ItemID = itemID

		// TODO: requires addition of UUID varchar(255) w/ unique index
		// TODO: reversal
		stmt, err := db.Prepare(`
			INSERT INTO check_ins (order_id, item_id, ticket_number, scanned_at, location, uuid, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)

		if err != nil {
			log.Println(err.Error())
		}
		defer stmt.Close()

		// TODO: ensure a row was created
		_, err = stmt.Exec(strconv.Itoa(scan.OrderID), strconv.Itoa(scan.ItemID), scan.TicketNumber, scan.Time, scan.Location, scan.UUID, timestamp, timestamp)
		if err != nil {
			log.Println(err.Error())
		}

		savedUUIDs = append(savedUUIDs, scan.UUID)
	}

	return
}
