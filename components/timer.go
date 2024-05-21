package components

import (
	"strconv"
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

type Timer struct {
	sync.Mutex
	workingDuration, breakingDuration time.Duration
	last                              time.Time
	working                           bool
	breaking                          bool
	Ratio                             float64
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
	t.breakingDuration += time.Duration(float64(t.workingDuration) / t.Ratio)
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
		m("hr", nil),
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
						"value": timefuncs.FormatDuration(time.Duration(float64(t.workingDuration) / t.Ratio)),
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
		m("tr", nil,
			m("td", nil, m("button#startWork.control", js.M{
				"config": mithril.RouteConfig,
				"onclick": func() {
					t.StartWork()
				},
				"disabled": t.working,
			}, s("Start working"))),

			m("td", nil,
				m("label[for='ratio']", nil,
					s("Break Ratio")),
				m("br", nil),
				m("select#ratio.center", js.M{
					"onchange": mithril.WithAttr("value", func(value string) {

						ratio, err := strconv.Atoi(value)
						if err == nil {
							t.Ratio = float64(ratio)
						} else {
							print("Failed to parse value", value)
						}

					})},
					m("option[value='1']", nil, s("1")),
					m("option[value='2']", nil, s("2")),
					m("option[value='3'][selected=true]", nil, s("3")),
					m("option[value='4']", nil, s("4")),
					m("option[value='5']", nil, s("5")),
				)),
			m("td", nil,
				m("button#startBreak.control", js.M{
					"config": mithril.RouteConfig,
					"onclick": func() {
						t.StartBreak()
					},
					"disabled": t.breaking,
				},
					s("Start Break")))),
		m("tr", nil,
			m("td[colspan=3]", nil,
				m("button#bigBreak.button.center", js.M{
					"config": mithril.RouteConfig,
					"onclick": func() {
						t.BigBreak()
					},
				},
					s("Big Break--reset all clocks"))),
		))
}
