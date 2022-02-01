package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

const size = 3

type Vector struct {
	X float64
	Y float64
}

type Point struct {
	N int64
	X float64
	Y float64
	Z float64
}

func dot(a Point, b Vector) float64 {
	res := a.X*b.X + a.Y*b.Y
	return res
}

func lerp(a float64, b float64, t int) float64 {
	res := a + (b-a)*float64(t)
	return res
}

func smoothstep(t int) int {
	return t * t * (3 - 2*t)
}

func unitVector() Vector {
	phi := 2 * math.Pi * rand.Float64()
	res := Vector{math.Cos(phi), math.Sin(phi)}
	return res
}

func setz(i int, j int, model [size][size]Point, vector [size][size]Vector) float64 {
	var nexti int
	if i == 0 {
		nexti = size - 1
	} else {
		nexti = i - 1
	}
	var nextj int
	if j == 0 {
		nextj = size - 1
	} else {
		nextj = j - 1
	}
	d1 := dot(model[i][j], vector[i][j])
	d2 := dot(model[nexti][j], vector[nexti][j])
	d3 := dot(model[i][nextj], vector[i][nextj])
	d4 := dot(model[nexti][nextj], vector[nexti][nextj])
	step1 := lerp(d1, d2, smoothstep(i))
	step2 := lerp(d3, d4, smoothstep(i))
	return lerp(step1, step2, smoothstep(j))

}

func main() {
	rand.Seed(time.Now().UnixNano())
	field := [size][size]Vector{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			field[i][j] = unitVector()
		}
	}
	model := [size][size]Point{}
	var n int64 = 1
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			model[i][j] = Point{
				n,
				float64(i),
				0.0,
				float64(j),
			}
			n++
		}
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			model[i][j].Y = setz(i, j, model, field)
		}
	}
	fmt.Println(model)
	f, err := os.OpenFile("test.obj", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			temp := model[i][j]
			f.WriteString(fmt.Sprintf("v %f %f %f\n", temp.X, temp.Y, temp.Z))
		}
	}
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-1; j++ {
			f.WriteString(fmt.Sprintf("f %d %d %d %d\n", model[i][j].N, model[i][j+1].N, model[i+1][j+1].N, model[i+1][j].N))
		}
	}

}
