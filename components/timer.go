package components

import (
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"

	"thirdtimer/timefuncs"
)

var m = moria.M

type s = moria.S

var (
	fps30 = time.Tick(time.Second / 30)
)

const RATIO = 3

type Timer struct {
	sync.Mutex
	workingDuration, breakingDuration time.Duration
	last                              time.Time
	working                           bool
	breaking                          bool
}

func (t *Timer) Controller() moria.Controller {
	return t
}

func (t *Timer) StartWork() {
	t.Lock()
	defer t.Unlock()
	t.last = time.Now()
	t.working = true
	t.breaking = false

	go func() {
		for t.breaking {
		}
		for t.working {
			<-fps30
			now := time.Now()
			t.workingDuration += now.Sub(t.last)
			t.last = now
			mithril.Redraw(false)
		}
	}()
}

// StartBreak stops the current work timer, adds the amount you've just earned
// to your breaking timer, and then switches the mode.
func (t *Timer) StartBreak() {
	t.Lock()
	defer t.Unlock()
	t.last = time.Now()
	t.working = false
	t.breaking = true
	t.breakingDuration += t.workingDuration / RATIO
	t.workingDuration = 0

	go func() {
		for t.working {
		}
		for t.breaking {
			<-fps30
			now := time.Now()
			t.breakingDuration -= now.Sub(t.last)
			t.last = now
			mithril.Redraw(false)
		}
	}()
}

func (t *Timer) BigBreak() {
	t.Lock()
	defer t.Unlock()
	t.last = time.Now()
	t.working = false
	t.breaking = false
	t.workingDuration = 0
	t.breakingDuration = 0
}

func (*Timer) View(ctrl moria.Controller) moria.View {
	t := ctrl.(*Timer)

	maybeRed := js.M{}
	if 0 < t.breakingDuration {
		maybeRed["style"] = "color:darkred;"
	}

	return m("div#wrapper", nil,
		m("h1", nil, s("Third Timer")),
		m("div", nil, s("See "),
			m("a", js.M{"href": "https://www.lesswrong.com/posts/RWu8eZqbwgB9zaerh/third-time-a-better-way-to-work"}, s("Third Time: A better way to work")),
			s(" for more information an how to use this timer."),
		),
		m("table#meetingTicker", nil,
			m("tr", nil,
				m("td", nil,
					m("label.copy[for='timeWorked']", nil, s("Time worked this stint"))),
				m("td.rightCell", nil,
					m("input#timeWorked", js.M{
						"value": timefuncs.FormatDuration(t.workingDuration),
					}))),

			m("tr", nil,
				m("td", nil,
					m("label.copy[for='timeEarned']", nil, s("Time earned this stint"))),
				m("td.rightCell", nil,
					m("input#timeEarned", js.M{
						"value": timefuncs.FormatDuration(t.workingDuration / RATIO),
					}))),

			m("tr", nil,
				m("td", nil,
					m("label.copy[for='timeEarned']", nil, s("Break time available"))),
				m("td.rightCell", nil,
					m("input#timeEarned", js.M{
						"value": timefuncs.FormatDuration(t.breakingDuration),
						"style": maybeRed["style"],
					}))),
		),
		m("button#startWork.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				t.StartWork()
			},
		}, s("Start working")),
		m("button#startBreak.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				t.StartBreak()
			},
		},
			s("Start Break")),
		m("button#bigBreak.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				t.BigBreak()
			},
		},
			s("Big Break--reset all clocks")),
	)
}
