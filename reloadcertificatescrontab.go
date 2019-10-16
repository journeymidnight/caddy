package caddy

import (
	"bytes"
	"fmt"
	"github.com/journeymidnight/yig-front-caddy/caddyfile"
	"github.com/journeymidnight/yig-front-caddy/telemetry"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
)

var IsCronTab bool

func CertificateCron(cdyfile Input) error {
	var CronTab []caddyfile.Token
	var times []string
	var Time string
	stypeName := cdyfile.ServerType()

	sblocks, err := loadServerBlocks(stypeName, cdyfile.Path(), bytes.NewReader(cdyfile.Body()))
	if err != nil {
		return err
	}

	for _, sb := range sblocks {
		for _, key := range sb.Keys {
			if key == ":443" {
				IsCronTab = true
				CronTab = sb.Tokens[caddyfile.TIMEFLAG]
			}
		}
	}
	
	if len(CronTab) == 0 {
		return nil
	}
	
	for n, value := range CronTab {
		if n == 0 && value.Text != caddyfile.TIMEFLAG {
			return fmt.Errorf("Wrong time parameter!")
		}
		if n > 6 {
			return fmt.Errorf("Wrong time parameter!")
		}
		times = append(times, value.Text)
	}
	Time = strings.Join(times[1:], " ")

	if IsCronTab {
		fmt.Println("Set the time for reloading the certificate to :", Time)
		go func() {
			c := cron.New()
			c.AddFunc(Time, ReloadGoroutine)
			c.Start()
		}()
		fmt.Println("Enable the method of periodically loading certificates from the database!")
	}
	return nil
}

func ReloadGoroutine() {
	go telemetry.AppendUnique("sigtrap", "SIGUSR1")

	// Start with the existing Caddyfile
	caddyfileToUse, inst, err := getCurrentCaddyfile()
	if err != nil {
		log.Printf("[ERROR] SIGUSR1: %v", err)
		return
	}
	if loaderUsed.loader == nil {
		// This also should never happen
		log.Println("[ERROR] SIGUSR1: no Caddyfile loader with which to reload Caddyfile")
		return
	}

	// Load the updated Caddyfile
	newCaddyfile, err := loaderUsed.loader.Load(inst.serverType)
	if err != nil {
		log.Printf("[ERROR] SIGUSR1: loading updated Caddyfile: %v", err)
		return
	}
	if newCaddyfile != nil {
		caddyfileToUse = newCaddyfile
	}

	// Backup old event hooks
	oldEventHooks := cloneEventHooks()

	// Purge the old event hooks
	purgeEventHooks()

	// Kick off the restart; our work is done
	EmitEvent(InstanceRestartEvent, nil)
	_, err = inst.Restart(caddyfileToUse)
	if err != nil {
		restoreEventHooks(oldEventHooks)
		log.Printf("[ERROR] SIGUSR1: %v", err)
		return
	}
}
