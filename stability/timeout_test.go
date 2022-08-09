package stability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTimeout(t *testing.T) {
	tests := []struct {
		name     string
		setUpCtx func() context.Context
		f        SlowFunction[string]
		wantErr  bool
		want     string
	}{
		{
			"1. timeout exited",
			func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond)
				return ctx
			},
			func(s string) (string, error) {
				time.Sleep(time.Microsecond)
				return "some string", nil
			},
			true,
			"",
		},
		{
			"2. got error",
			context.Background,
			func(s string) (string, error) {
				return "", errors.New("some error")
			},
			true,
			"",
		},
		{
			"3. success",
			context.Background,
			func(s string) (string, error) {
				return "success", nil
			},
			false,
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Timeout(tt.f)(tt.setUpCtx(), "some str")
			if (err != nil) != tt.wantErr {
				t.Errorf("Tomeout() got err: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Timeout() mismatch {+want; -got}: %s", diff)
			}
		})
	}
}
