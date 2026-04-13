// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package errorutil

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumererror"
)

func TestGetHTTPStatusCodeFromError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "retryable error maps to service unavailable",
			err:  errors.New("retryable"),
			want: http.StatusServiceUnavailable,
		},
		{
			name: "permanent error maps to bad request",
			err:  consumererror.NewPermanent(errors.New("permanent")),
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, GetHTTPStatusCodeFromError(tt.err))
		})
	}
}

func TestHTTPErrorWritesStatusAndMessage(t *testing.T) {
	t.Parallel()

	err := consumererror.NewPermanent(errors.New("bad payload"))
	recorder := httptest.NewRecorder()

	HTTPError(recorder, err)

	resp := recorder.Result()
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	require.NoError(t, readErr)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, "Permanent error: bad payload\n", string(body))
}

func TestHTTPErrorWithNilErrorDoesNothing(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()

	HTTPError(recorder, nil)

	resp := recorder.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Empty(t, string(body))
}
