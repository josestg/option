package option

import "testing"

type Record struct{ ID int }

func TestOption_String(t *testing.T) {
	expectEq(t, None[Record]().String(), "None")
	expectEq(t, Some(42).String(), "Some(42)")
	expectEq(t, Some(Record{ID: 1}).String(), "Some(option.Record{ID:1})")
}

func TestOption_Value(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		opt := Some(42)
		expectEq(t, opt.Value(), 42)
	})

	t.Run("absent", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic, got nil")
			}
		}()

		None[int]().Value()
	})
}

func TestOption_Absent(t *testing.T) {
	expectEq(t, None[int]().Absent(), true)
	expectEq(t, Some(42).Absent(), false)
}

func TestOption_Present(t *testing.T) {
	expectEq(t, None[int]().Present(), false)
	expectEq(t, Some(42).Present(), true)
}

func TestOption_ValueOr(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		opt := Some(42)
		expectEq(t, opt.ValueOr(0), 42)
	})

	t.Run("absent", func(t *testing.T) {
		expectEq(t, None[int]().ValueOr(42), 42)
	})
}

func TestOption_ValueOrBy(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		opt := Some(42)
		expectEq(t, opt.ValueOrBy(func() int { return 0 }), 42)
	})

	t.Run("absent", func(t *testing.T) {
		expectEq(t, None[int]().ValueOrBy(func() int { return 42 }), 42)
	})
}

func TestOption_Alt(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		opt := Some(42)
		expectEq(t, opt.Alt(None[int]()), opt)
	})

	t.Run("absent", func(t *testing.T) {
		opt := None[int]()
		alt := Some(42)
		expectEq(t, opt.Alt(alt).Value(), alt.Value())
	})
}

func TestOption_AltBy(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		opt := Some(42)
		expectEq(t, opt.AltBy(func() Option[int] { return None[int]() }), opt)
	})

	t.Run("absent", func(t *testing.T) {
		opt := None[int]()
		alt := Some(42)
		expectEq(t, opt.AltBy(func() Option[int] { return alt }).Value(), alt.Value())
	})
}

func TestOption_MarshalJSON(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		b, err := None[int]().MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, string(b), `{"kind":"None","value":0}`)
	})

	t.Run("some", func(t *testing.T) {
		b, err := Some(42).MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, string(b), `{"kind":"Some","value":42}`)
	})

	t.Run("some record", func(t *testing.T) {
		b, err := Some(Record{ID: 1}).MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, string(b), `{"kind":"Some","value":{"ID":1}}`)
	})

	t.Run("none record", func(t *testing.T) {
		b, err := None[Record]().MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, string(b), `{"kind":"None","value":{"ID":0}}`)
	})
}

func TestOption_UnmarshalJSON(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var opt Option[int]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"None","value":0}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, opt.Absent(), true)
	})

	t.Run("some", func(t *testing.T) {
		var opt Option[int]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"Some","value":42}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectEq(t, opt.Present(), true)
		expectEq(t, opt.Value(), 42)
	})

	t.Run("none record", func(t *testing.T) {
		var opt Option[Record]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"None","value":{"ID":0}}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectEq(t, opt.Absent(), true)
	})

	t.Run("some record", func(t *testing.T) {
		var opt Option[Record]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"Some","value":{"ID":1}}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectEq(t, opt.Present(), true)
		expectEq(t, opt.Value(), Record{ID: 1})
	})

	t.Run("unexpected kind", func(t *testing.T) {
		var opt Option[int]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"Maybe","value":42}`)); err == nil {
			t.Errorf("expected error, got nil")
		}

		expectEq(t, opt.Absent(), true)
	})

	t.Run("unexpected value", func(t *testing.T) {
		var opt Option[Record]
		if err := opt.UnmarshalJSON([]byte(`{"kind":"Some","value":42}`)); err == nil {
			t.Errorf("expected error, got nil")
		}

		expectEq(t, opt.Absent(), true)
	})
}

func expectEq[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}
