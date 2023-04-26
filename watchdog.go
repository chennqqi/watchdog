// Package watchdog implements control of hardware watchdog devices.
package watchdog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type WDT_OPTIONS uint32

const (
	WDIOS_DISABLECARD WDT_OPTIONS = 0x0001 /* Turn off the watchdog timer */
	WDIOS_ENABLECARD  WDT_OPTIONS = 0x0002 /* Turn on the watchdog timer */
	WDIOS_TEMPPANIC   WDT_OPTIONS = 0x0004 /* Kernel panic on temperature trip */

	WDIOC_SETTIMEOUT    = 0xc0045706
	WDIOC_SETPRETIMEOUT = 0xc0045708

	WDIOF_OVERHEAT   = 0x0001 /* Reset due to CPU overheat */
	WDIOF_FANFAULT   = 0x0002 /* Fan failed */
	WDIOF_EXTERN1    = 0x0004 /* External relay 1 */
	WDIOF_EXTERN2    = 0x0008 /* External relay 2 */
	WDIOF_POWERUNDER = 0x0010 /* Power bad/power fault */
	WDIOF_CARDRESET  = 0x0020 /* Card previously reset the CPU */
	WDIOF_POWEROVER  = 0x0040 /* Power over voltage */
	WDIOF_SETTIMEOUT = 0x0080 /* Set timeout (in seconds) */
	WDIOF_MAGICCLOSE = 0x0100 /* Supports magic close char */
	WDIOF_PRETIMEOUT = 0x0200 /* Pretimeout (in seconds), get/set */
	WDIOF_ALARMONLY  = 0x0400 /* Watchdog triggers a managemen */
)

// errNotImplemented is a sentinel which indicates package watchdog does not
// support this OS.
var errNotImplemented = fmt.Errorf("watchdog: not implemented on %s: %w", runtime.GOOS, os.ErrNotExist)

// A Device is a hardware watchdog device which can be pinged to keep the system
// from rebooting once the device has been opened.
type Device struct {
	// Identity is the name of the watchdog driver.
	Identity string

	f *os.File
}

type WatchdogInfo struct {
	Options  uint32
	Version  uint32
	Identity [32]uint8
}

func (t *WatchdogInfo) GetIdentity() string {
	return strings.TrimRight(string(t.Identity[:]), "\x00")
}

// Open opens the primary watchdog device on this system ("/dev/watchdog" on
// Linux, TBD on other platforms). If the device is not found, an error
// compatible with os.ErrNotExist will be returned.
//
// Once a Device is opened, you must call Ping repeatedly to keep the system
// from being rebooted. Call Close to disarm the watchdog device.
func Open(device_path string) (*Device, error) { return open(device_path) }

// Timeout returns the configured timeout of the watchdog device.
func (d *Device) BootStatus() (int, error) { return d.getBootStatus() }

// Timeout returns the pre timeout of the watchdog device.
func (d *Device) PreTimeout() (time.Duration, error) { return d.getPretimeout() }

// Timeout returns the pre timeout of the watchdog device.
func (d *Device) SetPreTimeout(t time.Duration) error { return d.setPretimeout(t) }

// Timeout returns the status of the watchdog device.
func (d *Device) Status() (int, error) { return d.getStatus() }

// Timeout returns the support of the watchdog device.
func (d *Device) Support() (*WatchdogInfo, error) { return d.getSupport() }

// All watchdog drivers are required return more information about the system,
// some do temperature, fan and power level monitoring, some can tell you the reason
// for the last reboot of the system. The GETSUPPORT ioctl is available to ask what
// the device can do:
func (d *Device) Temp() (int, error) { return d.getTemp() }

// Timeout returns the left timeout of the watchdog device.
func (d *Device) Timeleft() (time.Duration, error) { return d.getTimeLeft() }

// Timeout returns the configured timeout of the watchdog device.
func (d *Device) Timeout() (time.Duration, error) { return d.getTimeout() }

// Timeout returns the configured timeout of the watchdog device.
func (d *Device) SetTimeout(t time.Duration) error { return d.setTimeout(t) }

func (d *Device) KeepAlive() error { return d.keepAlive() }

func (d *Device) SetOptions(opts WDT_OPTIONS) error { return d.setOptions(opts) }

// Close closes the handle to the watchdog device and attempts to gracefully
// disarm the device, so no further Ping calls are required to keep the system
// from rebooting.
func (d *Device) Close() error { return d.close() }
