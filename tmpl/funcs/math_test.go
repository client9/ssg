package funcs

import "testing"

func TestToFloat64(t *testing.T) {
	cases := []struct {
		in   any
		want float64
		ok   bool
	}{
		{int(3), 3.0, true},
		{int64(-7), -7.0, true},
		{float32(1.5), 1.5, true},
		{float64(2.5), 2.5, true},
		{uint(4), 4.0, true},
		{"3.14", 3.14, true},
		{"bad", 0, false},
		{true, 0, false},
	}
	for _, c := range cases {
		got, err := toFloat64(c.in)
		if c.ok && err != nil {
			t.Errorf("toFloat64(%v): unexpected error: %v", c.in, err)
		}
		if !c.ok && err == nil {
			t.Errorf("toFloat64(%v): expected error, got %v", c.in, got)
		}
		if c.ok && got != c.want {
			t.Errorf("toFloat64(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestMathOps(t *testing.T) {
	fm := FuncMap()

	call2 := func(name string, a, b any) float64 {
		t.Helper()
		fn := fm[name].(func(any, any) (float64, error))
		v, err := fn(a, b)
		if err != nil {
			t.Fatalf("%s(%v, %v): %v", name, a, b, err)
		}
		return v
	}
	call1 := func(name string, a any) float64 {
		t.Helper()
		fn := fm[name].(func(any) (float64, error))
		v, err := fn(a)
		if err != nil {
			t.Fatalf("%s(%v): %v", name, a, err)
		}
		return v
	}

	if got := call2("add", 3, 4); got != 7 {
		t.Errorf("add: got %v", got)
	}
	if got := call2("sub", 10, 3); got != 7 {
		t.Errorf("sub: got %v", got)
	}
	if got := call2("mul", 3, 4); got != 12 {
		t.Errorf("mul: got %v", got)
	}
	if got := call2("div", 7, 2); got != 3.5 {
		t.Errorf("div: got %v", got)
	}
	if got := call2("mod", 7, 3); got != 1 {
		t.Errorf("mod: got %v", got)
	}
	if got := call1("abs", -5); got != 5 {
		t.Errorf("abs: got %v", got)
	}
	if got := call1("ceil", 1.2); got != 2 {
		t.Errorf("ceil: got %v", got)
	}
	if got := call1("floor", 1.9); got != 1 {
		t.Errorf("floor: got %v", got)
	}
	if got := call1("round", 1.5); got != 2 {
		t.Errorf("round: got %v", got)
	}
}

func TestMathDiv_byZero(t *testing.T) {
	if _, err := mathDiv(1, 0); err == nil {
		t.Error("expected error for div by zero")
	}
}

func TestMathMod_byZero(t *testing.T) {
	if _, err := mathMod(5, 0); err == nil {
		t.Error("expected error for mod by zero")
	}
}

func TestMathOps_mixedTypes(t *testing.T) {
	got, err := mathDiv(float64(7), int(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.5 {
		t.Errorf("div(7.0, 2) = %v, want 3.5", got)
	}
}
