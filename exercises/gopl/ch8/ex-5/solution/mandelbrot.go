/*
Упражнение 8.5. Возьмите существующую последовательную программу, такую
как программа вычисления множества Мандельброта из раздела 3.3 или вычисления
трехмерной поверхности из раздела 3.2, и выполните ее главный цикл параллельно, с
использованием каналов. Насколько быстрее стала работать программа на многопро­
цессорной машине? Каково оптимальное количество используемых go-подпрограмм?
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"time"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	type res struct {
		px, py int
		z      complex128
	}
	worksNum := 4
	var ch = make(chan res, width)
	var d = make(chan struct{}, worksNum)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for w := 0; w < worksNum; w++ {
		go func() {
			for c := range ch {
				img.Set(c.px, c.py, mandelbrot(c.z))
			}
			d <- struct{}{}
		}()
	}
	start := time.Now()

	for py := 0; py < height; py++ {
		r := res{}
		r.py = py
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			r.px, r.z = px, z
			ch <- r
		}
	}
	close(ch)
	data, err := os.Create("data.png")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	for w := 0; w < worksNum; w++ {
		<-d
	}
	png.Encode(data, img) // NOTE: ignoring errors
	fmt.Println(time.Since(start))
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
