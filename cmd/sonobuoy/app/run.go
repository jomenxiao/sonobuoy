/*
Copyright 2018 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"os"

	ops "github.com/heptio/sonobuoy/cmd/sonobuoy/app/operations"
	"github.com/heptio/sonobuoy/pkg/errlog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var runopts ops.RunConfig

func init() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Submits a sonobuoy run",
		Run:   submitSonobuoyRun,
	}
	cmd.PersistentFlags().StringVar(
		&runopts.Mode, "mode", "Conformance",
		"TBD: Update description on different run modes (quick|conformance|extended)",
	)
	// TODO: We should expose FOCUS and other options with sane defaults
	RootCmd.AddCommand(cmd)
}

func submitSonobuoyRun(cmd *cobra.Command, args []string) {
	code := 0
	if err := ops.Run(runopts); err != nil {
		errlog.LogError(errors.Wrap(err, "error attempting to run sonobuoy"))
		code = 1
	}
	os.Exit(code)
}
