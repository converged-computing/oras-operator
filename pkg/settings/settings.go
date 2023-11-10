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

// getDefaultSettings ensures that we don't update a global variable
func getDefaultSettings() map[string]OrasCacheSetting {
	return map[string]OrasCacheSetting{

		// Files are expected to be copied to/from here
		"input-path":  {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},
		"output-path": {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},
		"output-pipe": {Required: false, NonEmpty: true, Value: defaults.DefaultMissing},

		// Input and output container URIs for input/output artifacts
		// The input URI can be a listing (pulling from one or more dependnecy steps)
		"input-uri":  {Required: false, NonEmpty: true, Value: defaults.DefaultMissing, Listing: true, Values: []string{}},
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
}

type OrasCacheSetting struct {
	Required bool

	// If required (and provided) it cannot be empty
	NonEmpty bool

	// Listings
	Listing bool
	Values  []string

	Value string
}

// Oras Cache Settings are parsed from annotations
type Settings map[string]OrasCacheSetting

type ParsedSetting struct {
	IsList bool
	Field  string
}

// parseAnnotation handles parsing an ORAS operator annotation field into the field
// We also determine if it is a list.
func parseAnnotation(key string) *ParsedSetting {

	// If there are two slashes, this indicates a list item
	var field string

	// Indicates that this is a list value
	listValue := false

	// Pattern is to allow > 1 with <prefix>/<field>_<count>
	if strings.Contains(key, "_") {

		// We don't currently use the last identifier but could
		parts := strings.SplitN(key, "/", 2)
		field = parts[1]
		parts = strings.SplitN(field, "_", 2)
		field = parts[0]
		listValue = true
	} else {
		parts := strings.SplitN(key, "/", 2)
		field = parts[1]
	}

	return &ParsedSetting{
		IsList: listValue,
		Field:  field,
	}
}

type OrasCacheSettings struct {
	MarkedForOras   bool
	Settings        Settings
	DefaultSettings map[string]OrasCacheSetting
}

// Get a named setting
func (s *OrasCacheSettings) Get(name string) string {
	setting, ok := s.Settings[name]

	// If not defined, return NA
	if !ok {
		return s.getDefaultSetting(name)
	}
	return setting.Value
}

func (s *OrasCacheSettings) GetList(name string) []string {
	setting, ok := s.Settings[name]

	// If not defined, return NA
	if !ok {
		return s.getDefaultListSetting(name)
	}
	return setting.Values
}

// getDefaultSetting gets the default setting, if exists.
func (s *OrasCacheSettings) getDefaultSetting(name string) string {

	setting, ok := s.DefaultSettings[name]

	// If we know the setting, return the default value
	if ok {
		return setting.Value
	}
	// Otherwise we have no idea.
	return ""
}

// getDefaultSetting gets the default setting, if exists.
func (s *OrasCacheSettings) getDefaultListSetting(name string) []string {

	setting, ok := s.DefaultSettings[name]

	// If we know the setting, return the default value
	if ok {
		return setting.Values
	}
	// Otherwise we have no idea.
	return []string{}
}

// PrintSettings print all settings if debug mode is on
func (s *OrasCacheSettings) PrintSettings() {
	for name, setting := range s.Settings {
		if setting.Listing && setting.Value != "" {
			logger.Infof("🌟️ %s: %s", name, setting.Value)
		} else if setting.Listing {
			logger.Infof("🌟️ %s: %s", name, strings.Join(setting.Values, "\n"))
		} else {
			logger.Infof("🌟️ %s: %s", name, setting.Value)
		}
	}
}

func (s *OrasCacheSettings) Validate() bool {

	// Show the user the settings (for debugging)
	logger.Info(s.Settings)
	for key, ds := range s.DefaultSettings {

		// Retrieve the default, no go if required
		setting, ok := s.Settings[key]

		// If we don't have it, and it's required but a default provided
		if !ok && ds.Required && !ds.Listing && ds.Value != "" {
			s.Settings[key] = ds
			continue
		}

		// Same, but a listing
		if !ok && ds.Required && ds.Listing && len(ds.Values) == 0 {
			s.Settings[key] = ds
			continue
		}

		if !ok && ds.Required {
			logger.Warnf("The %s/%s annotation is required", defaults.OrasCachePrefix, key)
		}

		// Continue (ignore) if setting is not required
		if !ok {
			continue
		}
		if ds.NonEmpty && !ds.Listing && setting.Value == "" {
			logger.Warnf("The %s/%s is empty, and cannot be.", defaults.OrasCachePrefix, key)
			return false
		}
		// A listing can also be provided as a single value
		if ds.NonEmpty && ds.Listing && len(setting.Values) == 0 && setting.Value == "" {
			logger.Warnf("The %s/%s is empty, and cannot be.", defaults.OrasCachePrefix, key)
			return false
		}
	}

	// We don't require either of input or output paths
	return true
}

// NewOrasCacheSettings creates new settings
func NewOrasCacheSettings(annotations map[string]string) *OrasCacheSettings {

	// Create settings with defaults
	wrapper := OrasCacheSettings{}
	settings := Settings{}
	defaultSettings := getDefaultSettings()

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

			parsed := parseAnnotation(key)
			if parsed.Field == "debug" && value == "true" {
				debug = true
			}

			defaultSetting, ok := defaultSettings[parsed.Field]
			if !ok {
				logger.Warnf("Setting %s is not known the the oras operator.", key)
				continue
			}
			// Don't add the value if an empty string
			// TODO double check this does not alter default settings
			wrapper.MarkedForOras = true

			// Add a regular or list value, and update default Settings so we retrieve next time
			if parsed.IsList {
				logger.Infof("Setting %s is a list.", key)
				defaultSetting.Values = append(defaultSetting.Values, value)
				defaultSettings[parsed.Field] = defaultSetting

				// It's a list, but we were given a value (e.g., input-uri)
			} else {
				logger.Infof("Setting %s is not a list.", key)
				defaultSetting.Value = value
			}
			settings[parsed.Field] = defaultSetting
		}
	}
	wrapper.Settings = settings
	wrapper.DefaultSettings = defaultSettings
	if debug {
		wrapper.PrintSettings()
	}
	return &wrapper
}
