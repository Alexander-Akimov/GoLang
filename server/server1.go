package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	}

	http.HandleFunc("/", handler1) // Каждый запрос вызывет обработчик

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Обработчик возврщает компонент пути из URL запроса.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(os.Stderr, "URL.Path = %q\n", r.URL.Path)
}

var palette = []color.Color{color.Black, color.White,
	color.RGBA{0x00, 0xff, 0x00, 0xff}}

const (
	whiteIndex = 0 // Первый цвет палитры
	blackIndex = 1 // Следующий цвет палитры
	greenIndex = 2 // Зеленый - третий цвет
)

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // Количество полных колебаний
		res     = 0.001 // Угловое разрешение
		size    = 100   // Канва изображения охватывает [size..+size]
		nframes = 64    // Количество кадров анимации
		delay   = 8     // задержка между кадрами (единица - 10 мс)
	)
	rand.Seed(time.Now().UTC().UnixNano())
	freq := rand.Float64() * 3.0 // Отоносительная частота колебаний y
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // Разность фаз
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // Примечание: игнорируем ошибки

}
