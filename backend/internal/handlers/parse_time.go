package handlers

import (
	"time"

	"github.com/araddon/dateparse"
)

type ParserTime time.Time

func (ft *ParserTime) UnmarshalJSON(b []byte) error {

	s := string(b)
	s = s[1 : len(s)-1]

	t, err := dateparse.ParseAny(s)
	if err != nil {
		return err
	}

	*ft = ParserTime(t)

	return nil
}
