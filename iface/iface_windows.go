package iface

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/driver"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
	"net"
)

// Create Creates a new Wireguard interface, sets a given IP and brings it up.
func (w *WGIface) Create() error {

	WintunStaticRequestedGUID, _ := windows.GenerateGUID()
	adapter, err := driver.CreateAdapter(w.Name, "WireGuard", &WintunStaticRequestedGUID)
	if err != nil {
		err = fmt.Errorf("error creating adapter: %w", err)
		return err
	}
	w.Interface = adapter
	luid := adapter.LUID()
	err = adapter.SetAdapterState(driver.AdapterStateUp)
	if err != nil {
		return err
	}
	state, _ := luid.GUID()
	log.Debugln("device guid: ", state.String())
	return w.assignAddr(luid)
}

// assignAddr Adds IP address to the tunnel interface and network route based on the range provided
func (w *WGIface) assignAddr(luid winipcfg.LUID) error {

	log.Debugf("adding address %s to interface: %s", w.Address.IP, w.Name)
	err := luid.SetIPAddresses([]net.IPNet{{w.Address.IP, w.Address.Network.Mask}})
	if err != nil {
		return err
	}

	return nil
}
