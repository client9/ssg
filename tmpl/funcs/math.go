package funcs

import (
	"fmt"
	"math"
	"strconv"
	"text/template"
)

func mathFuncMap() template.FuncMap {
	return template.FuncMap{
		"add":   func(a, b any) (float64, error) { return applyOp(a, b, func(x, y float64) float64 { return x + y }) },
		"sub":   func(a, b any) (float64, error) { return applyOp(a, b, func(x, y float64) float64 { return x - y }) },
		"mul":   func(a, b any) (float64, error) { return applyOp(a, b, func(x, y float64) float64 { return x * y }) },
		"div":   mathDiv,
		"mod":   mathMod,
		"abs":   func(a any) (float64, error) { return applyFunc(a, math.Abs) },
		"ceil":  func(a any) (float64, error) { return applyFunc(a, math.Ceil) },
		"floor": func(a any) (float64, error) { return applyFunc(a, math.Floor) },
		"round": func(a any) (float64, error) { return applyFunc(a, math.Round) },
	}
}

func mathDiv(a, b any) (float64, error) {
	x, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	y, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	if y == 0 {
		return 0, fmt.Errorf("div: division by zero")
	}
	return x / y, nil
}

func mathMod(a, b any) (float64, error) {
	x, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	y, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	if y == 0 {
		return 0, fmt.Errorf("mod: division by zero")
	}
	return math.Mod(x, y), nil
}

func applyOp(a, b any, op func(float64, float64) float64) (float64, error) {
	x, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	y, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	return op(x, y), nil
}

func applyFunc(a any, fn func(float64) float64) (float64, error) {
	x, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	return fn(x), nil
}

// toFloat64 converts any numeric type or numeric string to float64.
func toFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case string:
		return strconv.ParseFloat(n, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}
