package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/v3/data/inventory"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

var errBinaryNotFound = errors.New("could not find apache binary")

func getBinPath(binPath string) (string, error) {
	var paths []string
	if binPath != "" {
		paths = append(paths, binPath)
	}
	paths = append(paths, "/usr/sbin/httpd", "/usr/sbin/apache2ctl")
	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			log.Debug("Using apache binary %q", path)
			return path, nil
		}
		log.Debug("Probing for apache binary at %q failed: %v", path, err)
	}
	return "", errBinaryNotFound
}

// setInventory executes system command in order to retrieve required inventory data and calls functions which parse the result.
// It returns a map of inventory data
func setInventory(inventory *inventory.Inventory, configBinaryPath string) error {
	commandPath, err := getBinPath(configBinaryPath)
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
