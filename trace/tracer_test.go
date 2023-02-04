package trace

import (
  "bytes"
  "testing"
)

// TestNew testuje działanie pakietu tracer
func TestNew(t *testing.T) {

  var buf bytes.Buffer
  tracer := New(&buf)

  if tracer == nil {
    t.Error("Funkcja New nie może zwracać nil!")
  }
  tracer.Trace("Witamy w pakiecie trace.")
  if buf.String() != "Witamy w pakiecie trace.\n" {
    t.Errorf("Metoda Trace nie powinna wyświetlić '%s'.", buf.String())
  }

}

func TestOff(t *testing.T) {
  silentTracer := Off()
  silentTracer.Trace("cośtam")
}
