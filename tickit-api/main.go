// Tickit Check-in Web API
package main

import (
	"code.google.com/p/gcfg"
	"encoding/json"
	"flag"
	"github.com/hoisie/web"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var Tickets = map[string]Ticket{}
var config Config
var storeID int

// respond to CORS OPTIONS check
func setHeaders(ctx *web.Context) {
	// TODO: should this be done with nginx
	// TODO: check for existence in config.Hosts.AllowedOrigins
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
	ctx.SetHeader("Access-Control-Allow-Credentials", "true", true)
	ctx.SetHeader("Access-Control-Allow-Methods", "OPTIONS, GET, POST", true)
	ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control, X-Tickit-Store, X-Tickit-Key", true)
	ctx.SetHeader("Content-type", "application/json; charset=utf-8", true)
}

func postScans(ctx *web.Context) {
	setHeaders(ctx)
	if authorize(ctx) == false {
		return
	}

	var scans []Scan

	response := RecordAPIResponse{Success: false}

	body, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		response.Message = "Error reading body"
		out, _ := json.Marshal(response)
		ctx.Abort(400, string(out))
		return
	}

	err = json.Unmarshal(body, &scans)
	if err != nil {
		response.Message = "Error decoding JSON"
		out, _ := json.Marshal(response)
		ctx.Abort(400, string(out))
		return
	}

	uuids, err := SaveScans(scans)
	if err != nil {
		response.Message = err.Error()
		out, _ := json.Marshal(response)
		ctx.Abort(500, string(out))
		return
	}

	response.Success = true
	response.SavedUUIDs = uuids
	out, _ := json.Marshal(response)
	// TODOO ctx.ResponseWriter.WriteHeader(201)
	ctx.Write(out)
}

// Parse common query params.
//
// Expects `item_ids` to be a comma-separated string.
//
// Expects `since` to be an ISOString in UTC (e.g. `2006-01-02T15:04:05.99Z`)
// which can be generated in JavaScript with `(new Date).toISOString()`.
func getParams(ctx *web.Context) (itemIDs []int, since time.Time) {
	itemIDs = make([]int, 0)

	raw := ctx.Request.FormValue("item_ids")

	strs := strings.Split(raw, ",")

	for _, str := range strs {
		itemID, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		itemIDs = append(itemIDs, itemID)
	}

	raw = ctx.Request.FormValue("since")

	since, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		ctx.Server.Logger.Printf("error parsing time (%s), using `0`", raw)
	}

	return
}

func getManifests(ctx *web.Context) {
	setHeaders(ctx)
	if authorize(ctx) == false {
		return
	}

	response := ManifestsAPIResponse{Last: time.Now().UTC(), Success: true}

	if response.Success {
		manifests, err := LoadManifests()

		if err != nil {
			panic(err.Error())
		}

		response.Manifests = manifests
	}

	out, err := json.Marshal(response)

	if err != nil {
		panic(err.Error())
	}

	ctx.Write(out)
}

func getTickets(ctx *web.Context) {
	setHeaders(ctx)
	if authorize(ctx) == false {
		return
	}

	itemIDs, since := getParams(ctx)

	response := TicketsAPIResponse{Since: since, Success: true}

	if len(itemIDs) < 1 {
		response.Message = "item_ids param was missing or invalid: must be a comma separated list of IDs"
		response.Success = false
	}

	if response.Success {
		response.Tickets, response.Last = LoadTickets(itemIDs, since)
	}

	out, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	ctx.Write(out)
}

func getScans(ctx *web.Context) {
	setHeaders(ctx)
	if authorize(ctx) == false {
		return
	}

	itemIDs, since := getParams(ctx)

	response := ScansAPIResponse{Since: since, Success: true}

	if len(itemIDs) < 1 {
		response.Message = "item_ids param was missing or invalid: must be a comma separated list of IDs"
		response.Success = false
	}

	if since.IsZero() {
		since = time.Unix(0, 0)
	}

	if response.Success {
		response.Scans, response.Last = LoadScans(itemIDs, since)
	}

	out, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}
	ctx.Write(out)
}

func authorize(ctx *web.Context) (ok bool) {
	ok = true
	if ctx.Request.Method == "OPTIONS" {
		return
	}
	store := ctx.Request.Header.Get("X-Tickit-Store")
	key := ctx.Request.Header.Get("X-Tickit-Key")

	if store == "" || key == "" {
		ctx.Server.Logger.Printf("missing auth headers", store, key)
		ok = false
	}

	if ok {
		// TODO: assigning a global
		storeID, err := strconv.Atoi(store)
		if err != nil {
			ctx.Server.Logger.Printf("error parsing X-Tickit-Store (%s)", store)
			ok = false
		} else {
			ok = StoreExists(storeID, key)
		}
	}

	if ok == false {
		response := GenericAPIResponse{Success: false, Message: "What are you doing, Dave?"}
		out, _ := json.Marshal(response)
		ctx.WriteHeader(401)
		ctx.Write(out)
	}

	return ok
}

func main() {
	configPath := flag.String("config", "./checkin-api.conf", "path to the configuration file")
	flag.Parse()
	log.Printf("configPath: %+v", *configPath)
	err := gcfg.ReadFileInto(&config, *configPath)

	if err != nil {
		panic("Unable to load checkin-api.conf")
	}

	log.Printf("loaded config: %+v", config)

	go startSocketServer()
	web.Get("/manifests/?", getManifests)
	web.Get("/tickets/?", getTickets)
	web.Get("/scans/?", getScans)
	web.Post("/record/?", postScans)
	web.Match("OPTIONS", "/.*", setHeaders)
	web.Run(config.Servers.Http)
}
