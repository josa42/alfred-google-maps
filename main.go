package main

// Package is called aw
import (
	"fmt"
	"net/url"
	"strings"

	"github.com/deanishe/awgo"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	query := strings.Split(strings.Join(wf.Args(), " "), ">")
	icon := &aw.Icon{Value: "icon/icon-main.png"}

	start := ""
	end := ""
	waypoints := []string{}

	for idx, item := range query {
		item = strings.Trim(item, " ")

		if idx == 0 {
			start = item
		} else if idx == len(query)-1 {
			end = item
		} else {
			waypoints = append(waypoints, item)
		}

	}

	if end == "" && len(waypoints) == 0 {
		wf.NewItem(start).
			Subtitle("Search on Google Maps").
			Icon(icon).
			Arg(fmt.Sprintf(
				"https://www.google.com/maps/search/?api=1&query=%s",
				url.QueryEscape(start),
			)).
			Valid(true)

	} else {
		startLabel := start
		if startLabel == "" {
			startLabel = "Current Location"
		}

		viaLabel := ""
		if len(waypoints) > 0 {
			viaLabel = fmt.Sprintf(" (via %s)", strings.Join(waypoints, ", "))
		}

		wf.NewItem(fmt.Sprintf("From %s to %s%s", startLabel, end, viaLabel)).
			Subtitle("Directions on Google Maps").
			Icon(icon).
			Arg(fmt.Sprintf(
				"https://www.google.com/maps/dir/?api=1&origin=%s&destination=%s&waypoints=%s",
				url.QueryEscape(start),
				url.QueryEscape(end),
				url.QueryEscape(strings.Join(waypoints, "|")),
			)).
			Valid(true)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
