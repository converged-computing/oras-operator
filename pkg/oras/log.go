/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package oras

import (
	"log"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	handle, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = handle.Sugar()
	defer handle.Sync()
}
