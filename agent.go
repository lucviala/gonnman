package connman

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
)

type Agent struct {
	Service    string
	Path       dbus.ObjectPath
	Interface  string
	Name       string
	Passphrase string
}

func NewAgent(ssid, psk string) *Agent {
	agent := &Agent{
		Service:    "net.gonnman",
		Path:       "/net/connman/Agent",
		Interface:  "net.connman.Agent",
		Name:       ssid,
		Passphrase: psk,
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	reply, err := conn.RequestName(agent.Service, dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "Name already taken")
		return nil
	}

	conn.Export(agent, agent.Path, agent.Interface)
	return agent
}

func (a *Agent) Destroy() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	reply, err := conn.ReleaseName(a.Service)
	if err != nil {
		return err
	}
	if reply != dbus.ReleaseNameReplyReleased {
		return fmt.Errorf("Could not release the name\n")
	}

	conn.Export(nil, a.Path, a.Interface)
	return nil
}

func (a *Agent) RequestInput(service dbus.ObjectPath, rq map[string]dbus.Variant) (map[string]dbus.Variant, *dbus.Error) {
	var in map[string]dbus.Variant
	if a.Name != "" {
		in = map[string]dbus.Variant{
			"Name":       dbus.MakeVariant(a.Name),
			"Passphrase": dbus.MakeVariant(a.Passphrase),
		}
	} else {
		in = map[string]dbus.Variant{
			"Passphrase": dbus.MakeVariant(a.Passphrase),
		}
	}
	return in, nil
}

func (a *Agent) ReportError(service dbus.ObjectPath, err string) *dbus.Error {
	log.Printf("%s: %s\n", service, err)
	return nil
}
