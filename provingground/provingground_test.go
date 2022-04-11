package provingground

import (
	"math"
	"math/rand"
	v "rlbb/lib/vector"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Circle struct {
	pos    v.V2
	speed  v.V2
	radius float64

	vspeed float64
}

type Missile struct {
	pos    v.V2
	speed  v.V2
	vspeed float64
}

// http://opensteer.sourceforge.net/
// https://stackoverflow.com/questions/17729737/how-to-find-an-intercept-on-a-moving-target/17749335#17749335
//https://stackoverflow.com/questions/51851120/2d-target-intercept-algorithm
// P + t*U = Q + t*V
// U = V + (Q-P)/t
// s*s = |U|
// s*s = V.V + 2*(Q-P).V/t + (Q-P).(Q-P)/(t*t)
//(s*s - V.V)*t*t - 2*(Q-P).V*t - (Q-P).(Q-P) = 0
//ax2+bx+c = 0
//delta = b2-4ac
// t1 = (-b + sqr(delta))/2a
// t1 = (-b - sqr(delta))/2a
//a=(s*s - V.V)
//b=-2*(Q-P).V
//c= -(Q-P).(Q-P)
//delta = 4*(Q-P).V*(Q-P).V+4*(s*s - V.V)*(Q.P).(Q-P)
//qp = Q-p
//delta = 4*qp.V*qp.V+4*(U.len2()-V.V)*qp.qp

func rnd() float64 { return rand.Float64() }
func Test1(t *testing.T) {
	rl.SetTraceLog(rl.LogNone)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)
	rl.InitWindow(1440, 720, "test")
	//missile, speed V2
	rl.SetTargetFPS(60)

	// V := circle.speed.MulA(circle.vspeed)
	// U := missile.speed.MulA(missile.vspeed)
	// P := missile.pos
	// Q := circle.center

	// s2 := U.Len2()
	// QP := Q.Sub(P)
	// a := s2 - V.NormDot(V)
	// b := -2 * QP.NormDot(V)
	// c := -(QP.NormDot(QP))
	// delta := b*b - 4*a*c
	// var t1, t2, time float64
	// fmt.Println("delta=", delta)
	// if delta >= 0 {
	// 	t1 = (-b + math.Sqrt(delta)) / (2 * a)
	// 	t2 = (-b - math.Sqrt(delta)) / (2 * a)
	// 	fmt.Println("Solutions: t1=", t1, " t2=", t2)
	// }
	// if t1 > 0 {
	// 	time = t1
	// }
	// if t2 > 0 {
	// 	time = math.Max(time, t2)
	// }

	var circle Circle
	var missile Missile
	var intersect v.V2
	//var cp v.V2
	var time float64

	for !rl.WindowShouldClose() {

		if rl.IsKeyDown(rl.KeySpace) {

			circle = Circle{v.V2{X: 1200 * rnd(), Y: 100 * rnd()}, v.V2{rnd(), rnd()}.Norm(), 10 + rnd()*30, 0.5 + 2*rnd()}
			missile = Missile{v.V2{750, 400}, v.V2{}, 4}
			//time = findPlayerIntercept(circle.pos, circle.speed.MulA(circle.vspeed), missile.pos, missile.vspeed)
			time = findPlayerIntercept2(circle, missile)
			intersect = circle.pos.Add(circle.speed.MulA(time * circle.vspeed))
			missile.speed = intersect.Sub(missile.pos).Norm()

			//			cp = circle.pos.Add(circle.speed.MulA(time * circle.vspeed))
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
 
		rl.DrawCircleLines(int32(circle.pos.X), int32(circle.pos.Y), float32(circle.radius), rl.Yellow)

		rl.DrawCircle(int32(missile.pos.X), int32(missile.pos.Y), 10, rl.White)
		rl.EndDrawing()

		dist := missile.pos.Sub(circle.pos).Len()
		if time > 0 && dist >= circle.radius {
			circle.pos = circle.pos.Add(circle.speed.MulA(circle.vspeed))
			missile.pos = missile.pos.Add(missile.speed.MulA(missile.vspeed))
		}
	}

	t.Error("OK")
}

func findPlayerIntercept(circlePos v.V2, circleVel v.V2, missilePos v.V2, missileVel float64) float64 {
	// calculate the speeds
	v := missileVel * missileVel
	u := math.Sqrt(circleVel.X*circleVel.X + circleVel.Y*circleVel.Y)
	// calculate square distance
	c := (missilePos.X-circlePos.X)*(missilePos.X-circlePos.X) +
		(missilePos.Y-circlePos.Y)*(missilePos.Y-circlePos.Y)

	// calculate first two quadratic coefficients
	a := v*v - u*u
	b := circleVel.X*(missilePos.X-circlePos.X) + circleVel.Y*(missilePos.Y-circlePos.Y)

	// collision time
	t := -1.0 // invalid value

	// if speeds are equal
	if math.Abs(a) < 0.00005 { // some small number, e.g. 1e-5f
		t = c / (2.0 * b)
	} else {
		// discriminant
		b /= a
		d := b*b + c/a

		// real roots exist
		if d > 0.0 {
			// if single root
			if math.Abs(d) < 0.00005 {
				t = b / a
			} else {
				// how many positive roots?
				e := math.Sqrt(d)
				if math.Abs(b) < e {
					t = b + e
				} else if b > 0.0 {
					t = b - e
				}
			}
		}
	}
	return t

}
func findPlayerIntercept2(c Circle, m Missile) float64 {

	t := float64(0)
	for t < 1440 {

		cp := c.pos.Add(c.speed.MulA(t * c.vspeed))
		dist := cp.Sub(m.pos).Len()

		if dist <= m.vspeed*t {
			return float64(t)
		}
		t = t + 1
	}
	return -1
}
