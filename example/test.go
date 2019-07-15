package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jiyu93/gopid"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

// Usage :
// Click http://127.0.0.1:28080/test?p=50&i=2.5&d=1.5 and adjust p、i、d
func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/test", paint)
	err := http.ListenAndServe(":28080", nil)
	if err != nil {
		panic(err)
	}
}

func paint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	kp, _ := strconv.ParseFloat(r.Form["p"][0], 10)
	ki, _ := strconv.ParseFloat(r.Form["i"][0], 10)
	kd, _ := strconv.ParseFloat(r.Form["d"][0], 10)
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "gopid Testing"
	p.X.Label.Text = "time"
	p.Y.Label.Text = "speed"
	points := runPID(kp, ki, kd)
	err = plotutil.AddLinePoints(
		p,
		"target", points[0],
		"real speed", points[1],
	)
	if err != nil {
		panic(err)
	}
	c, err := p.WriterTo(1200, 500, "jpg")
	if err != nil {
		panic(err)
	}
	c.WriteTo(w)
}

// runPID : stimulate an object runing from 0km/s to 80km/s, in uniformly accelerated motion
// m=1000kg
// dt=1s, 180s in total
// F=ma --> a=F/m
// v1=v0+a*dt
// output = F
func runPID(p, i, d float64) []plotter.XYs {
	var xAxis = 180
	var targetSpeed float64 = 80 * 1000
	var m float64 = 1000
	var v float64

	points := make([]plotter.XYs, 2)
	points[0] = make(plotter.XYs, xAxis)
	points[1] = make(plotter.XYs, xAxis)

	c := gopid.NewPID(p, i, d, targetSpeed)

	for i := 0; i < xAxis; i++ {
		// X axis
		points[0][i].X = float64(i)
		points[1][i].X = float64(i)
		// target line
		points[0][i].Y = targetSpeed
		// real speed line
		F := c.CalcIncPID(v)
		if F > 5000000 {
			F = 5000000
		}
		a := F / m
		v = v + a*1
		points[1][i].Y = v

		// make some error
		v += rand.Float64() * 10 * c.Kp
		v -= rand.Float64() * 10 * c.Kp

		// change target
		if i == 60 {
			targetSpeed = 40 * 1000
			c.TargetValue = targetSpeed
		}
		if i == 120 {
			targetSpeed = 120 * 1000
			c.TargetValue = targetSpeed
		}
	}
	return points
}
