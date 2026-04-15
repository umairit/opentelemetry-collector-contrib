// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sampling // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor/internal/sampling"

import "go.opentelemetry.io/collector/featuregate"

var _ = featuregate.GlobalRegistry().MustRegister(
	"processor.tailsamplingprocessor.disableinvertdecisions",
	featuregate.StageStable,
	featuregate.WithRegisterDescription("When enabled, sampling policy 'invert_match' will result in a SAMPLED or NOT SAMPLED decision instead of INVERT SAMPLED or INVERT NOT SAMPLED."),
	featuregate.WithRegisterToVersion("v0.152.0"),
)
