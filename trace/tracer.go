package trace

import (
  "fmt"
  "io"
)

// Tracer jest inrefejsem opisującym obiekt pozwalającym na 
// śledzenie zdarzeń zachodzących podczas wykonywania kodu.
type Tracer interface {
  Trace(...interface{})
}

// New creates a new Tracer that will write the output to
// the specified io.Writer.
func New(w io.Writer) Tracer {
  return &tracer{out: w}
}

// tracer jest implementacją interfejsu Tracer, 
// która zapisuje io.Writer.
type tracer struct {
  out io.Writer
}

// Trace zapisuje argumenty do implementacji io.Writer.
func (t *tracer) Trace(a ...interface{}) {
  fmt.Fprint(t.out, a...)
  fmt.Fprintln(t.out)
}

// nilTracer
type nilTracer struct{}

// Metoda Trace w tym przypadku nic nie robi.
func (t *nilTracer) Trace(a ...interface{}) {}

// Metoda Off tworzy implementację Tracer, która będzie ignorować
// wywołania metody Trace.
func Off() Tracer {
  return &nilTracer{}
}
