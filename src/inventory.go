package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/data/inventory"
	"io"
	"os"
	"os/exec"
	"strings"
)

func getBinPath() (string, error) {
	// Check first for RedHat
	binPath := "/usr/sbin/httpd"
	_, err := os.Stat(binPath)
	if err != nil {
		// If it isn't a RedHat, try with Debian
		binPath = "/usr/sbin/apache2ctl"
		_, derr := os.Stat(binPath)
		if derr != nil {
			return "", fmt.Errorf("it isn't possible to locate Apache executable")
		}
	}
	return binPath, nil
}

// setInventory executes system command in order to retrieve required inventory data and calls functions which parse the result.
// It returns a map of inventory data
func setInventory(inventory *inventory.Inventory) error {
	commandPath, err := getBinPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(commandPath, "-M")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching the inventory (modules). Message: %v", err.Error())
	}
	r := bytes.NewReader(output)
	err = getModules(bufio.NewReader(r), inventory)
	if err != nil {
		return err
	}

	cmd = exec.Command(commandPath, "-V")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching the inventory (version). Message: %v", err.Error())
	}
	r = bytes.NewReader(output)
	err = getVersion(bufio.NewReader(r), inventory)
	if err != nil {
		return err
	}

	if len(inventory.Items()) == 0 {
		return fmt.Errorf("empty result")
	}
	return nil
}

// getModules reads an Apache list of enabled modules and transforms its
// contents into a map that can be processed by NR agent.
// It appends a map of inventory data where the keys contain name of the module and values
// indicate that module is enabled.
func getModules(reader *bufio.Reader, i *inventory.Inventory) error {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(line, "_module") {
			splitedLine := strings.Split(line, "_module")
			moduleName := strings.TrimSpace(splitedLine[0])
			key := fmt.Sprintf("modules/%s", moduleName)
			err = i.SetItem(key, "value", "enabled")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getVersion reads an Apache list of compile settings and transforms its
// contents into a map that can be processed by NR agent.
// It appends a map of inventory data which indicates Apache Server version
func getVersion(reader *bufio.Reader, i *inventory.Inventory) error {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(line, "Server version") {
			splitedLine := strings.Split(line, ":")
			err = i.SetItem("version", "value", strings.TrimSpace(splitedLine[1]))
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
