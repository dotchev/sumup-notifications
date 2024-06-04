package itest

import (
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/require"
)

func Retry(t *testing.T, f func() error) {
	err := retry.Do(f,
		retry.Delay(100*time.Millisecond),
		retry.Attempts(50),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true))
	require.NoError(t, err)
}
