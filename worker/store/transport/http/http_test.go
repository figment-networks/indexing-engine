package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundRobin(t *testing.T) {
	t.Run("round robin next url", func(t *testing.T) {

		urls := []string{"url1", "url2", "url3", "url4"}
		rr := newRoundRobin(urls)

		// run several rounds
		for i := 1; i <= 3; i++ {
			for ii := 0; ii < len(urls); ii++ {
				require.Equal(t, urls[ii], rr.getNext())
			}
		}
	})
}
