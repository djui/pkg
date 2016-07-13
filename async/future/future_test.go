package future

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	task := func(...interface{}) (interface{}, error) {
		return "done", nil
	}

	f := New(task)
	v, err := f.Get()
	assertEqual(t, nil, err)
	assertEqual(t, "done", v)
}

func TestNewCompleted(t *testing.T) {
	f := NewCompleted("done", nil)
	v, err := f.Get()
	assertEqual(t, nil, err)
	assertEqual(t, "done", v)

	e := errors.New("error")

	f = NewCompleted(nil, e)
	v, err = f.Get()
	assertEqual(t, e, err)
	assertEqual(t, nil, v)
}

func TestIsDone(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}

	f := New(task)
	assertEqual(t, false, f.IsDone())
	assertEqual(t, false, f.IsDone())
	close(c)
	_, _ = f.Get()
	assertEqual(t, true, f.IsDone())
	assertEqual(t, true, f.IsDone())
}

func TestIsCanceled(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}

	f := New(task)
	assertEqual(t, false, f.IsCanceled())
	assertEqual(t, false, f.IsCanceled())
	f.Cancel()
	assertEqual(t, true, f.IsCanceled())
	assertEqual(t, true, f.IsCanceled())
}

func TestCancel(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}

	f := New(task)
	assertEqual(t, true, f.Cancel())
	assertEqual(t, false, f.Cancel())
	v, err := f.Get()
	assertEqual(t, ErrCanceled, err)
	assertEqual(t, nil, v)

	f = New(task)
	f.Complete("ignore")
	assertEqual(t, false, f.Cancel())
}

func TestComplete(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}

	f := New(task)
	assertEqual(t, true, f.Complete("override"))
	assertEqual(t, true, f.IsDone())
	assertEqual(t, false, f.IsCanceled())
	v, err := f.Get()
	assertEqual(t, nil, err)
	assertEqual(t, "override", v)

	assertEqual(t, false, f.Complete("override again"))
	assertEqual(t, true, f.IsDone())
	assertEqual(t, false, f.IsCanceled())
	v2, err2 := f.Get()
	assertEqual(t, nil, err2)
	assertEqual(t, v, v2)
}

func TestCompleteWithError(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}

	e := errors.New("error")

	f := New(task)
	assertEqual(t, true, f.CompleteWithError(e))
	assertEqual(t, true, f.IsDone())
	assertEqual(t, false, f.IsCanceled())
	v, err := f.Get()
	assertEqual(t, e, err)
	assertEqual(t, nil, v)

	assertEqual(t, false, f.CompleteWithError(fmt.Errorf("another error")))
	assertEqual(t, true, f.IsDone())
	assertEqual(t, false, f.IsCanceled())
	v2, err2 := f.Get()
	assertEqual(t, e, err2)
	assertEqual(t, nil, v2)
}

func TestGet(t *testing.T) {
	task := func(...interface{}) (interface{}, error) {
		return "done", nil
	}
	f := New(task)
	v, err := f.Get()
	assertEqual(t, nil, err)
	assertEqual(t, "done", v)
	v2, err2 := f.Get()
	assertEqual(t, err, err2)
	assertEqual(t, v, v2)
}

func TestGetNow(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}
	f := New(task)
	v, err := f.GetNow("fallback")
	assertEqual(t, nil, err)
	assertEqual(t, "fallback", v)
	close(c)
	_, _ = f.Get()
	v2, err2 := f.GetNow("another fallback")
	assertEqual(t, nil, err2)
	assertEqual(t, "done", v2)
}

func TestGetNowWithTimeout(t *testing.T) {
	c := make(chan struct{})
	task := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}
	f := New(task)
	v, err := f.GetWithTimeout(1 * time.Second)
	assertEqual(t, ErrTimeout, err)
	assertEqual(t, nil, v)
	close(c)
	v, err = f.GetWithTimeout(1 * time.Second)
	assertEqual(t, nil, err)
	assertEqual(t, "done", v)
}

func TestThenComposeSuccees(t *testing.T) {
	taskA := func(...interface{}) (interface{}, error) {
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + " and done", nil
	}
	taskC := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + " and over", nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCompose(taskB)
	f3 := f2.ThenCompose(taskC)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, nil, err2)
	assertEqual(t, "done and done", v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, false, f2.IsCanceled())

	v3, err3 := f3.Get()
	assertEqual(t, nil, err3)
	assertEqual(t, "done and done and over", v3)
	assertEqual(t, true, f3.IsDone())
	assertEqual(t, false, f3.IsCanceled())

}

func TestThenComposePartialFailure(t *testing.T) {
	e := errors.New("error")
	taskA := func(...interface{}) (interface{}, error) {
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return nil, e
	}
	taskC := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + " and over", nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCompose(taskB)
	f3 := f2.ThenCompose(taskC)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, e, err2)
	assertEqual(t, nil, v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, false, f2.IsCanceled())

	v3, err3 := f3.Get()
	assertEqual(t, e, err3)
	assertEqual(t, nil, v3)
	assertEqual(t, true, f3.IsDone())
	assertEqual(t, false, f3.IsCanceled())
}

func TestThenComposeCanceled(t *testing.T) {
	c := make(chan struct{})
	taskA := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + " and done", nil
	}
	taskC := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + " and over", nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCompose(taskB)
	f3 := f2.ThenCompose(taskC)

	f2.Cancel()
	close(c)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, ErrCanceled, err2)
	assertEqual(t, nil, v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, true, f2.IsCanceled())

	v3, err3 := f3.Get()
	assertEqual(t, ErrCanceled, err3)
	assertEqual(t, nil, v3)
	assertEqual(t, true, f3.IsDone())
	assertEqual(t, true, f3.IsCanceled())
}

func TestThenCombineSuccess(t *testing.T) {
	taskA := func(...interface{}) (interface{}, error) {
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return " and done", nil
	}
	combinator := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + args[1].(string), nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCombine(taskB, combinator)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, nil, err2)
	assertEqual(t, "done and done", v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, false, f2.IsCanceled())
}

func TestThenCombinePartialFailure(t *testing.T) {
	e := errors.New("error")
	taskA := func(...interface{}) (interface{}, error) {
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return nil, e
	}
	combinator := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + args[1].(string), nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCombine(taskB, combinator)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, e, err2)
	assertEqual(t, nil, v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, false, f2.IsCanceled())
}

func TestThenCombineCanceled(t *testing.T) {
	c := make(chan struct{})
	taskA := func(...interface{}) (interface{}, error) {
		<-c
		return "done", nil
	}
	taskB := func(args ...interface{}) (interface{}, error) {
		return " and done", nil
	}
	combinator := func(args ...interface{}) (interface{}, error) {
		return args[0].(string) + args[1].(string), nil
	}

	f1 := New(taskA)
	f2 := f1.ThenCombine(taskB, combinator)

	f2.Cancel()
	close(c)

	v1, err1 := f1.Get()
	assertEqual(t, nil, err1)
	assertEqual(t, "done", v1)
	assertEqual(t, true, f1.IsDone())
	assertEqual(t, false, f1.IsCanceled())

	v2, err2 := f2.Get()
	assertEqual(t, ErrCanceled, err2)
	assertEqual(t, nil, v2)
	assertEqual(t, true, f2.IsDone())
	assertEqual(t, true, f2.IsCanceled())
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		_, fn, line, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: %s != %s", filepath.Base(fn), line, expected, actual)
	}
}
