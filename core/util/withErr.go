package util

import "errors"

func WithErr[T any](perr *error) func(t T, err error) T {
	return func(t T, err error) T {
		if err != nil {
			if *perr != nil {
				*perr = errors.Join(*perr, err)
			} else {
				*perr = err
			}
		}
		return t
	}
}
