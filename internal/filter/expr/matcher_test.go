// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package expr

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type boolExprFunc[K any] func(context.Context, K) (bool, error)

func (f boolExprFunc[K]) Eval(ctx context.Context, tCtx K) (bool, error) {
	return f(ctx, tCtx)
}

type countingMatcher[K any] struct {
	result bool
	err    error
	calls  int
}

func (m *countingMatcher[K]) Eval(context.Context, K) (bool, error) {
	m.calls++
	return m.result, m.err
}

func TestAlwaysTrue(t *testing.T) {
	t.Parallel()

	matched, err := AlwaysTrue[int]().Eval(t.Context(), 42)
	require.NoError(t, err)
	require.True(t, matched)
}

func TestNot(t *testing.T) {
	t.Parallel()

	matchErr := errors.New("match failed")
	tests := []struct {
		name    string
		matcher BoolExpr[int]
		want    bool
		wantErr error
	}{
		{
			name: "inverts false result",
			matcher: boolExprFunc[int](func(context.Context, int) (bool, error) {
				return false, nil
			}),
			want: true,
		},
		{
			name: "preserves matcher error",
			matcher: boolExprFunc[int](func(context.Context, int) (bool, error) {
				return true, matchErr
			}),
			want:    false,
			wantErr: matchErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			matched, err := Not(tt.matcher).Eval(t.Context(), 1)
			require.Equal(t, tt.want, matched)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOrZeroMatchers(t *testing.T) {
	t.Parallel()

	require.Nil(t, Or[int]())
}

func TestOrOneMatcherReturnsOriginalMatcher(t *testing.T) {
	t.Parallel()

	matcher := &countingMatcher[int]{result: true}
	require.Same(t, matcher, Or[int](matcher))
}

func TestOrShortCircuitsOnFirstMatch(t *testing.T) {
	t.Parallel()

	first := &countingMatcher[int]{result: false}
	second := &countingMatcher[int]{result: true}
	third := &countingMatcher[int]{result: true}

	matched, err := Or[int](first, second, third).Eval(t.Context(), 1)
	require.NoError(t, err)
	require.True(t, matched)
	require.Equal(t, 1, first.calls)
	require.Equal(t, 1, second.calls)
	require.Zero(t, third.calls)
}

func TestOrReturnsFalseWhenNothingMatches(t *testing.T) {
	t.Parallel()

	first := &countingMatcher[int]{result: false}
	second := &countingMatcher[int]{result: false}

	matched, err := Or[int](first, second).Eval(t.Context(), 1)
	require.NoError(t, err)
	require.False(t, matched)
	require.Equal(t, 1, first.calls)
	require.Equal(t, 1, second.calls)
}

func TestOrStopsOnMatcherError(t *testing.T) {
	t.Parallel()

	matchErr := errors.New("boom")
	first := &countingMatcher[int]{err: matchErr}
	second := &countingMatcher[int]{result: true}

	matched, err := Or[int](first, second).Eval(t.Context(), 1)
	require.False(t, matched)
	require.ErrorIs(t, err, matchErr)
	require.Equal(t, 1, first.calls)
	require.Zero(t, second.calls)
}
