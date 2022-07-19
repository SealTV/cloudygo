package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SealTV/cloudygo/retry"
)

var count int

func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	}

	return "success", nil
}

func main() {
	r := retry.Retry(EmulateTransientError, 5, 2*time.Second)
	resp, err := r(context.Background())

	fmt.Println(resp, err)
}
