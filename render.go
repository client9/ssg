package ssg

import (
	"errors"
	"fmt"
)

// Pipeline is a named sequence of stages executed in order.
type Pipeline struct {
	name   string
	stages []Stage
}

// NewPipeline constructor
func NewPipeline(name string, stages ...Stage) Pipeline {
	return Pipeline{name: name, stages: stages}
}

func (p Pipeline) Run(ctx *Context, cfg ContentSourceConfig, in any) (any, error) {
	cur := in
	for _, s := range p.stages {
		next, err := s.Run(ctx, cfg, cur)
		if err != nil {
			return nil, fmt.Errorf("pipeline %q, stage %q failed: %w", p.name, s.Name(), err)
		}
		cur = next
	}
	return cur, nil
}

// Stage is a single named pipeline step with type-erased execution.
// Use Step to create a Stage from a typed function.
type Stage interface {
	Name() string
	Run(ctx *Context, cfg ContentSourceConfig, in any) (any, error)
}

// NewStage constructs a Stage from a name and a function.
func NewStage(name string, fn func(ctx *Context, cfg ContentSourceConfig, in any) (any, error)) Stage {
	return stageFunc{name: name, fn: fn}
}

type stageFunc struct {
	name string
	fn   func(*Context, ContentSourceConfig, any) (any, error)
}

func (s stageFunc) Name() string { return s.name }
func (s stageFunc) Run(ctx *Context, cfg ContentSourceConfig, in any) (any, error) {
	return s.fn(ctx, cfg, in)
}

// Step wraps a typed function into a Stage. The any type is confined here;
// stage functions use concrete types and are self-documenting.
//
// Pass-through stages (metadata-only, content unchanged) use [any, any]:
//
//	Step("set-output-file", func(ctx *Context, cfg ContentSourceConfig, in any) (any, error) { ... })
//
// Terminal sinks use struct{} as the output type:
//
//	Step("write-output", func(ctx *Context, cfg ContentSourceConfig, in []byte) (struct{}, error) { ... })
func Step[I, O any](name string, fn func(*Context, ContentSourceConfig, I) (O, error)) Stage {
	return stageFunc{
		name: name,
		fn: func(ctx *Context, cfg ContentSourceConfig, in any) (any, error) {
			typed, ok := in.(I)
			if !ok {
				return nil, fmt.Errorf("stage %q: got %T, want %T", name, in, *new(I))
			}
			return fn(ctx, cfg, typed)
		},
	}
}

// RunPipeline executes a Pipeline with the given input, returning the final output or any error encountered.
// The output is type-asserted to T, and a descriptive error is returned if the assertion fails.
// To just check for execution errors without caring about the output type, use T = any.
func RunPipeline[T any](ctx *Context, cfg ContentSourceConfig, p Pipeline, input any) (T, error) {
	out, err := p.Run(ctx, cfg, input)
	if err != nil {
		var zero T
		return zero, err // error already wrapped by Run
	}
	typed, ok := out.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("pipeline %q: final output got %T, want %T", p.name, out, zero)
	}
	return typed, nil
}

// Render is a Plugin that materializes all artifacts by running each one's
// Pipeline. Globals from ctx are merged into each artifact's Meta before the
// pipeline runs; page frontmatter wins on key collision.
func Render(ctx *Context, artifacts *[]Artifact) error {
	for i := range *artifacts {
		a := &(*artifacts)[i]
		if ctx != nil {
			for k, v := range ctx.Globals {
				if _, exists := a.Meta[k]; !exists {
					a.Meta[k] = v
				}
			}
		}
		source, _ := a.Meta["Content"].([]byte)
		if _, err := RunPipeline[any](ctx, a.Meta, a.Pipeline, source); err != nil {
			return fmt.Errorf("%s: %w", a.Meta.InputFile(), err)
		}
	}
	return nil
}

// Must wraps a (Stage, error) constructor result, panicking if err is non-nil.
// Use it to inline constructors that return errors in a pipeline slice literal.
//
//	ssg.Must(ssg.NewPageRender("layout", fns))
func Must(s Stage, err error) Stage {
	if err != nil {
		panic(err)
	}
	return s
}

// ---- Fan-out ---------------------------------------------------------------
//
// FanOut runs multiple branch Pipelines from a single input. All branches run
// regardless of individual failures — no short-circuit.
//
// Each branch receives a shallow clone of cfg so metadata mutations (e.g.
// SetOutputFile) in one branch do not affect the others. The content value
// (in) is shared across branches. Immutability cannot be enforced at this
// level because the type is erased, but a branch that needs its own copy
// can open with an explicit clone step:
//
//	var CloneBytes = ssg.Step("clone-bytes",
//	    func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
//	        return bytes.Clone(in), nil
//	    })
//
// Two usage patterns:
//
//   - Non-terminal: FanOut is followed by another Stage. On success it returns
//     (FanOutResult, nil) and the next stage receives the FanOutResult and can
//     inspect each branch's Name, Value, and Err freely.
//
//   - Terminal: FanOut is the last stage. On success the pipeline returns
//     cleanly. On failure it returns (nil, FanOutResult): the FanOutResult
//     implements the error interface so the full result — every branch outcome —
//     is preserved and recoverable with errors.As:
//
//     var far FanOutResult
//     if errors.As(err, &far) {
//     for _, b := range far.Branches {
//     fmt.Println(b.Name, b.Err)
//     }
//     }
//
// BranchResult holds the name, final output value, and error for one branch.
type BranchResult struct {
	Name  string
	Value any
	Err   error
}

// FanOutResult holds the outcome of every branch. It implements the error
// interface so it can be returned directly as an error when branches fail,
// making the full result set recoverable via errors.As even after the
// pipeline stops.
type FanOutResult struct {
	Branches []BranchResult
}

// Error implements the error interface. It is called only when at least one
// branch failed; the message is the join of all branch error strings.
func (r FanOutResult) Error() string {
	var errs []error
	for _, b := range r.Branches {
		if b.Err != nil {
			errs = append(errs, b.Err)
		}
	}
	return errors.Join(errs...).Error()
}

// Unwrap returns all branch errors, so errors.As and errors.Is traverse into
// individual branch failures.
func (r FanOutResult) Unwrap() []error {
	var errs []error
	for _, b := range r.Branches {
		if b.Err != nil {
			errs = append(errs, b.Err)
		}
	}
	return errs
}

// FanOut returns a Stage that runs each branch Pipeline with the same input.
// All branches run regardless of individual failures. On success it returns
// (FanOutResult, nil); on failure it returns (nil, FanOutResult) so the full
// result set is always recoverable via errors.As.
func FanOut(name string, branches ...Pipeline) Stage {
	return NewStage(name, func(ctx *Context, cfg ContentSourceConfig, in any) (any, error) {
		results := make([]BranchResult, 0, len(branches))

		for _, b := range branches {
			v, err := b.Run(ctx, cfg.Clone(), in)
			results = append(results, BranchResult{
				Name:  b.name,
				Value: v,
				Err:   err,
			})
		}

		result := FanOutResult{Branches: results}
		if len(result.Unwrap()) > 0 {
			return nil, result
		}
		return result, nil
	})
}
