package printing

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

// HeliosLog is log formatting for Styx
func HeliosLog(level, text string) {
	if level == "SYSTEM" {
		text = aurora.White("[*] ").String() + text
	}

	if level == "OPEN" {
		text = aurora.Green("[O] ").String() + text
	}

	if level == "CLOSED" {
		text = aurora.Red("[X] ").String() + text
	}

	if level == "ERROR" {
		text = aurora.Yellow("[E] ").String() + text
	}

	fmt.Println(text)
}
