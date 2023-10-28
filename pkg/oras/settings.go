/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	orasCachePrefix = "oras.converged-computing.github.io"
)

var (
	defaultSettings = map[string]OrasCacheSetting{
		"input-path":  {Required: false, NonEmpty: true},
		"output-path": {Required: false, NonEmpty: true},
		"identifier":  {Required: true, NonEmpty: true},
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

func (s *OrasCacheSettings) validate() bool {

	// Show the user the settings (for debugging)
	logger.Info(s.Settings)
	for key, defaultSetting := range defaultSettings {

		// Retrieve the default, no go if required
		setting, ok := s.Settings[key]
		if !ok && defaultSetting.Required {
			logger.Warnf("The %s/%s annotation is required", orasCachePrefix, key)
		}

		// Continue (ignore) if setting is not required
		if !ok {
			continue
		}
		if defaultSetting.NonEmpty && setting.Value == "" {
			logger.Warnf("The %s/%s is empty, and cannot be.", orasCachePrefix, key)
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
func NewOrasCacheSettings(pod *corev1.Pod) *OrasCacheSettings {

	// Create settings with defaults
	wrapper := OrasCacheSettings{}
	settings := Settings{}

	// Parse all annotations looking for oras cache prefix
	for key, value := range pod.Annotations {
		if strings.HasPrefix(key, orasCachePrefix) {

			// The annotation is required to be in format <identifier/field>
			if !strings.Contains(key, "/") {
				logger.Warnf("Provided key %s does not contain '/' to separate field, skipping.", key)
				continue
			}

			parts := strings.SplitN(key, "/", 2)
			field := parts[1]

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
	return &wrapper
}
