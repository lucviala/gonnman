package connman

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

type Technology struct {
	Path                dbus.ObjectPath
	Name                string
	Type                string
	Powered             bool
	Connected           bool
	Tethering           bool
	TetheringIdentifier string
	TetheringPassphrase string
}

func (t *Technology) Enable() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}
	return db.Set("Powered", true)
}

func (t *Technology) Disable() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}
	return db.Set("Powered", false)
}

func (t *Technology) Scan() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}

	_, err = db.Call("Scan")
	return err
}

func (t *Technology) EnableTethering(ssid string, psk string) error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}

	if len(ssid) > 0 {
		log.Printf("Setting up TetheringIdentifier: %v\n", ssid)
		db.Set("TetheringIdentifier", ssid)
		log.Printf("Setting up TetheringIdentifier: Successfully\n")
	}
	if len(psk) > 8 {
		log.Printf("Setting up TetheringPassphrase: %v\n", psk)
		db.Set("TetheringPassphrase", psk)
		log.Printf("Setting up TetheringPassphrase: Successfully\n")
	} else {
		return fmt.Errorf("Passphrase too short: %v", psk)
	}
	log.Printf("Enabling tethering: %v - %v\n", ssid, psk)
	return db.Set("Tethering", true)

}

func (t *Technology) DisableTethering() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	log.Println("Disabling tethering!")
	return db.Set("Tethering", false)
}
