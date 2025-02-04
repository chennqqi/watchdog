//go:build linux
// +build linux

package watchdog

import (
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// All system calls in this code are part of the Linux watchdog API. For
// reference, see:
// https://www.kernel.org/doc/html/latest/watchdog/watchdog-api.html.

func open(device_path string) (*Device, error) {
	// TODO(mdlayher): determine the significance of the "/dev/watchdogN" nodes
	// on Linux. It appears that my machine with only one device exposes both
	// "/dev/watchdog" and "/dev/watchdog0".
	//
	// According to Terin, /dev/watchdog is an alias for /dev/watchdog0 on
	// modern machines. It's possible there could be more than one device, so
	// we'll eventually want to support that.
	if device_path == "" {
		device_path = "/dev/watchdog"
	}
	f, err := os.OpenFile(device_path, os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	// Immediately fetch the device's information to return to the caller.
	info, err := unix.IoctlGetWatchdogInfo(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	// Immediately fetch the device's information to return to the caller.
	return &Device{

		// Clean up any trailing NULL bytes.
		Identity: strings.TrimRight(string(t.Identity[:]), "\x00"),

		f: f,
	}, nil
}

func (d *Device) setPretimeout(t time.Duration) error {
	return unix.IoctlSetInt(int(d.f.Fd()), WDIOC_SETPRETIMEOUT, int(t.Seconds()))
}

func (d *Device) setTimeout(t time.Duration) error {
	return unix.IoctlSetInt(int(d.f.Fd()), WDIOC_SETTIMEOUT, int(t.Seconds()))
}

func (d *Device) getBootStatus() (int, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETBOOTSTATUS)
	if err != nil {
		return 0, err
	}
	return s, err
}

func (d *Device) getPretimeout() (time.Duration, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETPRETIMEOUT)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return time.Duration(s) * time.Second, nil
}

func (d *Device) getStatus() (int, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETSTATUS)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return s, nil
}

func (d *Device) getSupport() (*WatchdogInfo, error) {
	// Immediately fetch the device's information to return to the caller.
	info, err := unix.IoctlGetWatchdogInfo(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	var dInfo WatchdogInfo
	dInfo.Options = info
	dInfo.Version = info.Version
	dInfo.Identity = append(dInfo.Identity, info.Identity)
	return &dInfo, nil
}

func (d *Device) getTemp() (int, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETTEMP)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return s, nil
}

func (d *Device) getTimeLeft() (time.Duration, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETTIMELEFT)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return time.Duration(s) * time.Second, nil
}

func (d *Device) getTimeLeft() (time.Duration, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETTIMELEFT)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return time.Duration(s) * time.Second, nil
}

func (d *Device) getTimeout() (time.Duration, error) {
	s, err := unix.IoctlGetInt(int(d.f.Fd()), unix.WDIOC_GETTIMEOUT)
	if err != nil {
		return 0, err
	}

	// The time value is always returned in seconds.
	return time.Duration(s) * time.Second, nil
}

func (d *Device) keepAlive() error { return unix.IoctlWatchdogKeepalive(int(d.f.Fd())) }

func (d *Device) setOptions(options WDT_OPTIONS) error {
	err := unix.IoctlSetInt(int(d.f.Fd()), unix.WDIOC_SETOPTIONS, int(options))
	if err != nil {
		return 0, err
	}
	return err
}

func (d *Device) close() error {
	// Attempt a Magic Close to disarm the watchdog device, since any call to
	// Close would be intentional and it's unlikely the user would want a system
	// reboot. Reference:
	// https://www.kernel.org/doc/html/latest/watchdog/watchdog-api.html#magic-close-feature
	if _, err := d.f.Write([]byte("V")); err != nil {
		// Make sure the file descriptor is closed even if Magic Close fails.
		_ = d.f.Close()
		return err
	}

	return d.f.Close()
}
