package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const defaultFmtStr = "{capacity}% - {usage}W - {H}h {M}m"

var (
	dischargeStr = flag.String("discharging", defaultFmtStr, "format string")
	chargeStr    = flag.String("charging", defaultFmtStr, "format string")
	fullStr      = flag.String("full", defaultFmtStr, "format string")
	waybarMode   = flag.Bool("waybar", false, "enable waybar mode")
)

func main() {
	flag.Parse()

	pc, err := NewPowerClient()

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create dbus conn: %v\n", err)
		os.Exit(-1)
	}

	bs, err := pc.GetBatteryStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get battery status: %v\n", err)
		os.Exit(-1)
	}

	fmtStr := *dischargeStr
	if bs.IsCharging && bs.Capacity > 99 {
		fmtStr = *fullStr
	} else if bs.IsCharging {
		fmtStr = *chargeStr
	}

	if *waybarMode {
		fmt.Fprintf(os.Stdout,
			`{"text":"%s","percentage": %1f,"class": "%s"}\n`,
			applyFmt(fmtStr, bs),
			bs.Capacity,
			strings.ReplaceAll(strings.ToLower(bs.State.String()), " ", "-"))

	} else {
		fmt.Fprintln(os.Stdout, applyFmt(fmtStr, bs))
	}
}

func applyFmt(str string, bs BatteryStatus) string {
	result := strings.ReplaceAll(str, "{capacity}", fmt.Sprintf("%.0f", bs.Capacity))

	eta := bs.TimeUntilEmpty
	if bs.IsCharging {
		eta = bs.TimeUntilFull
	}

	hours := "0"
	if eta.Hours() >= 1 {
		hours = fmt.Sprintf("%.0f", eta.Hours())
	}

	mins := "0"
	if eta.Minutes() >= 1 {
		mins = fmt.Sprintf("%1.0f", eta.Minutes())
	}

	result = strings.ReplaceAll(result, "{H}", hours)
	result = strings.ReplaceAll(result, "{M}", mins)
	result = strings.ReplaceAll(result, "{usage}", fmt.Sprintf("%.2f", bs.Usage))
	result = strings.ReplaceAll(result, "{state}", fmt.Sprintf("%s", bs.State))

	return result
}
