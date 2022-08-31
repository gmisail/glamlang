package util

import "math"

type Span struct {
	from int
	to   int
}

func CreateSpan(from int, to int) Span {
	return Span{from, to}
}

func (s Span) Combine(span Span) Span {
	from := int(math.Min(float64(s.from), float64(span.from)))
	to := int(math.Max(float64(s.to), float64(span.to)))

	return CreateSpan(from, to)
}
