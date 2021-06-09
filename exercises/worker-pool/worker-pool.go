package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var images = []string{"https://www.hdwallpapers.in/download/brown_black_white_cat_kitten_is_walking_on_sand_4k_hd_kitten-HD.jpg", "https://www.hdwallpapers.in/download/colorful_digital_art_swirl_shapes_abstraction_4k_hd_abstract-HD.jpg", "https://www.hdwallpapers.in/download/yellow_eyes_brown_white_cat_in_blur_background_hd_cat-HD.jpg", "https://www.hdwallpapers.in/download/beach_4k_2-HD.jpg", "https://www.hdwallpapers.in/download/red_plant_leaves_with_water_drops_4k_hd_nature-HD.jpg", "https://www.hdwallpapers.in/download/snow_covered_landscape_mountains_valley_trees_forest_under_blue_sky_4k_hd_nature-HD.jpg", "https://www.hdwallpapers.in/download/snow_covered_mountains_forest_slope_trees_4k_hd_nature-HD.jpg", "https://www.hdwallpapers.in/download/snow_covered_mountains_peaks_fog_under_purple_black_cloudy_sky_4k_hd_nature-HD.jpg", "https://www.hdwallpapers.in/download/white_mountains_valley_trees_under_blue_sky_4k_hd_nature-HD.jpg", "https://www.hdwallpapers.in/download/blue_orange_fractal_shapes_circles_glare_hd_abstract-HD.jpg"}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func worker(jobs <-chan string, done chan<- struct{}) {
	for j := range jobs {
		err := downloadFile(j, path.Base(j))
		if err != nil {
			fmt.Println(err)
			fmt.Println("couldn't did the job, url: %s" + j)
		}
		done <- struct{}{}
	}
}

func main() {
	start := time.Now()

	fmt.Println("number of images:", len(images))
	jobs := make(chan string, len(images))
	done := make(chan struct{}, len(images))

	var workersNum int
	if len(os.Args) == 2 {
		workersNum, _ = strconv.Atoi(os.Args[1])
	}
	if workersNum == 0 {
		workersNum = 2
	}

	for w := 0; w < workersNum; w++ {
		go worker(jobs, done)
	}

	for _, image := range images {
		jobs <- image
	}
	close(jobs)

	for range images {
		<-done
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
