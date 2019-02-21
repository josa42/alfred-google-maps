package main

// Package is called aw
import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

const (
	updateJobName = "checkForUpdate"
	repo          = "josa42/alfred-google-maps"
)

var (
	flagCheck     bool
	wf            *aw.Workflow
	iconAvailable = &aw.Icon{Value: "icon/update.png"}
)

func init() {
	wf = aw.New(update.GitHub(repo))

	flag.BoolVar(&flagCheck, "check", false, "Check for a new version")
}

func main() {
	wf.Run(run)
}

func run() {
	wf.Args()
	flag.Parse()

	if flagCheck {
		runCheck()
		return
	}

	runTriggerCheck()

	query := strings.Trim(strings.Join(wf.Args(), " "), " ")
	queryParts := strings.Split(query, ">")
	icon := &aw.Icon{Value: "icon.png"}

	start := ""
	end := ""
	waypoints := []string{}

	for idx, item := range queryParts {
		item = strings.Trim(item, " ")

		if idx == 0 {
			start = item
		} else if idx == len(queryParts)-1 {
			end = item
		} else {
			waypoints = append(waypoints, item)
		}

	}

	if query == "" {

		wf.NewItem("Google Maps").
			Subtitle("Usage: \"gm [location]\" or \"gm [origin] > [target]\"").
			Icon(icon).
			Valid(false)

		if wf.UpdateAvailable() {
			wf.Configure(aw.SuppressUIDs(true))

			wf.NewItem("Update available!").
				Subtitle("↩ to install").
				Autocomplete("workflow:update").
				Valid(false).
				Icon(iconAvailable)
		}
	} else if end == "" && len(waypoints) == 0 {
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

func runCheck() {
	wf.Configure(aw.TextErrors(true))
	log.Println("Checking for updates...")
	if err := wf.CheckForUpdate(); err != nil {
		wf.FatalError(err)
	}
}

func runTriggerCheck() {
	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}
}
