//go:build !linux
// +build !linux

package watchdog

import (
	"time"
)

func open(device_path string) (*Device, error) { return nil, errNotImplemented }

func (d *Device) setPretimeout(t time.Duration) error   { return errNotImplemented }
func (d *Device) setTimeout(t time.Duration) error      { return errNotImplemented }
func (d *Device) getBootStatus() (int, error)           { return 0, errNotImplemented }
func (d *Device) getPretimeout() (time.Duration, error) { return 0, errNotImplemented }
func (d *Device) getStatus() (int, error)               { return 0, errNotImplemented }
func (d *Device) getSupport() (*WatchdogInfo, error)    { return nil, errNotImplemented }
func (d *Device) getTemp() (int, error)                 { return 0, errNotImplemented }
func (d *Device) getTimeLeft() (time.Duration, error)   { return 0, errNotImplemented }
func (d *Device) getTimeout() (time.Duration, error)    { return 0, errNotImplemented }
func (d *Device) keepAlive() error                      { return errNotImplemented }
func (d *Device) setOptions(options WDT_OPTIONS) error  { return errNotImplemented }
func (*Device) close() error                            { return errNotImplemented }
