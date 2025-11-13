package model

import "github.com/go-faster/errors"

var (
	ErrOrderNotFound = errors.New("Order is not found")
	ErrCancelOrder   = errors.New("Error cancel order")
	ErrUpdateOrder   = errors.New("Error update order")
)
