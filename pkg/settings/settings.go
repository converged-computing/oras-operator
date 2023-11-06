/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package settings

import (
	"strings"

	"github.com/converged-computing/oras-operator/pkg/defaults"
)

var (
	defaultSettings = map[string]OrasCacheSetting{

		// Files are expected to be copied to/from here
		"input-path":  {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},
		"output-path": {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},
		"output-pipe": {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},

		// Input and output container URIs for input/output artifacts
		"input-uri":  {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},
		"output-uri": {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},

		// The name of the sidecar orchestrator
		"oras-cache": {Required: true, NonEmpty: true},

		// Debug mode to print / show all settings
		"debug": {Required: false, NonEmpty: true, Value: "false"},

		// The container with oras to run for the service
		"oras-container": {Required: true, Value: defaults.OrasBaseImage},

		// The name(s) of the launcher containers
		"container": {Required: false, NonEmpty: true},

		// Entrypoint custom script to wget
		"entrypoint":      {Required: false, NonEmpty: true, Value: defaults.ApplicationEntrypoint},
		"oras-entrypoint": {Required: false, NonEmpty: true, Value: defaults.OrasEntrypoint},
	}
)

type OrasCacheSetting struct {
	Required bool

	// If required (and provided) it cannot be empty
	NonEmpty bool
	Value    string
}

// Oras Cache Settings are parsed from annotations
type Settings map[string]OrasCacheSetting

type OrasCacheSettings struct {
	MarkedForOras bool
	Settings      Settings
}

// Get a named setting
func (s *OrasCacheSettings) Get(name string) string {
	setting, ok := s.Settings[name]

	// If not defined, return NA
	if !ok {
		return getDefaultSetting(name)
	}
	return setting.Value
}

// getDefaultSetting gets the default setting, if exists.
func getDefaultSetting(name string) string {

	setting, ok := defaultSettings[name]

	// If we know the setting, return the default value
	if ok {
		return setting.Value
	}
	// Otherwise we have no idea.
	return ""
}

// PrintSettings print all settings if debug mode is on
func (s *OrasCacheSettings) PrintSettings() {
	for name, setting := range s.Settings {
		logger.Infof("🌟️ %s: %s", name, setting.Value)
	}
}

func (s *OrasCacheSettings) Validate() bool {

	// Show the user the settings (for debugging)
	logger.Info(s.Settings)
	for key, defaultSetting := range defaultSettings {

		// Retrieve the default, no go if required
		setting, ok := s.Settings[key]

		// If we don't have it, and it's required but a default provided
		if !ok && defaultSetting.Required && defaultSetting.Value != "" {
			s.Settings[key] = defaultSetting
			continue
		}

		if !ok && defaultSetting.Required {
			logger.Warnf("The %s/%s annotation is required", defaults.OrasCachePrefix, key)
		}

		// Continue (ignore) if setting is not required
		if !ok {
			continue
		}
		if defaultSetting.NonEmpty && setting.Value == "" {
			logger.Warnf("The %s/%s is empty, and cannot be.", defaults.OrasCachePrefix, key)
			return false
		}
	}

	// One of input or output must be defined
	_, inputOk := s.Settings["input-path"]
	_, outputOk := s.Settings["output-path"]

	if !inputOk && !outputOk {
		logger.Warn("One of input-path or output-path is required.")
		return false
	}
	return true
}

// NewOrasCacheSettings creates new settings
func NewOrasCacheSettings(annotations map[string]string) *OrasCacheSettings {

	// Create settings with defaults
	wrapper := OrasCacheSettings{}
	settings := Settings{}

	// Do we have debug mode on?
	debug := false

	// Parse all annotations looking for oras cache prefix
	for key, value := range annotations {
		if strings.HasPrefix(key, defaults.OrasCachePrefix) {

			// The annotation is required to be in format <identifier/field>
			if !strings.Contains(key, "/") {
				logger.Warnf("Provided key %s does not contain '/' to separate field, skipping.", key)
				continue
			}

			parts := strings.SplitN(key, "/", 2)
			field := parts[1]
			if field == "debug" && value == "true" {
				debug = true
			}

			defaultSetting, ok := defaultSettings[field]
			if !ok {
				logger.Warnf("Setting %s is not known the the oras operator.", key)
				continue
			}
			// Don't add the value if an empty string
			// TODO double check this does not alter default settings
			wrapper.MarkedForOras = true
			defaultSetting.Value = value
			settings[field] = defaultSetting
		}
	}
	wrapper.Settings = settings
	if debug {
		wrapper.PrintSettings()
	}
	return &wrapper
}
