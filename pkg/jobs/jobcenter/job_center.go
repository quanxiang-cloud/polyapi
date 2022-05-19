package jobcenter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

var inst = &jobCenter{}

// RegRunner add a runner task to list
func RegRunner(order int, name string, desc string, runner Runner) error {
	return inst.RegRunner(order, name, desc, runner)
}

// ShowList show regstered task list
func ShowList() error {
	return inst.ShowList()
}

// Run start the jobs
func Run(db *gorm.DB) error {
	return inst.Run(db)
}

// Runner is a db job runner
type Runner interface {
	Run(db *gorm.DB) (string, error)
}

// RunnerFunc is a func that impliment Runner
type RunnerFunc func(db *gorm.DB) (string, error)

// Run is adaptor for Runner
func (f RunnerFunc) Run(db *gorm.DB) (string, error) {
	return f(db)
}

type job struct {
	order  int
	name   string
	desc   string
	runner Runner
}
type jobCenter struct {
	jobs   []*job
	sorted bool
}

func (c *jobCenter) sort() {
	if !c.sorted {
		sort.Slice(c.jobs, func(i, j int) bool { return c.jobs[i].order < c.jobs[j].order })
		c.sorted = true
	}
}

func (c *jobCenter) RegRunner(order int, name string, desc string, runner Runner) error {
	c.jobs = append(c.jobs, &job{
		name:   name,
		order:  order,
		desc:   desc,
		runner: runner,
	})
	return nil
}

func (c *jobCenter) ShowList() error {
	c.sort()
	fmt.Printf("==list of %d jobs:\n", len(c.jobs))
	fmt.Printf("%-3s %-20s %s\n", "", "NAME", "DESC")
	for i, v := range c.jobs {
		fmt.Printf("%-3d %-20s %s\n", i+1, v.name, v.desc)
	}
	return nil
}

func (c *jobCenter) Run(db *gorm.DB) error {
	c.sort()
	count := len(c.jobs)
	t := time.Now()
	last := t
	var report = strings.Builder{}
	if s := fmt.Sprintf("==%s running %d tasks\n", Now(t, nil), count); true {
		report.WriteString("\n")
		report.WriteString(s)
		fmt.Print(s)
	}
	for i, v := range c.jobs {
		fmt.Printf("====%s running %d/%d tasks %s", Now(t, nil), i+1, count, v.name)
		s, err := v.runner.Run(db)
		if err != nil {
			return fmt.Errorf("**%s:%s", v.name, err.Error())
		}
		if s := fmt.Sprintf("====%s finish %d/%d tasks %s, %s\n", Now(t, &last), i+1, count, v.name, s); true {
			report.WriteString(s)
			fmt.Print(s)
		}

	}
	if s := fmt.Sprintf("==%s finish %d tasks\n", Now(t, nil), count); true {
		report.WriteString(s)
		fmt.Print(report.String())
	}
	return nil
}

// Now get timestamp and duration
func Now(then time.Time, last *time.Time) string {
	t := time.Now()
	dur := ""
	if d := t.Sub(then) / time.Second * time.Second; d > 0 {
		dur = fmt.Sprintf("(%s)", d)
	}
	cost := ""
	if last != nil {
		elapse := t.Sub(*last) / time.Second * time.Second
		cost = fmt.Sprintf("(%s)", elapse)
		*last = t
	}

	return fmt.Sprintf("%s%s%s", t.Format("2006-01-02T15:04:05"), dur, cost)
}
