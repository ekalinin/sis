package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var (
	port       = flag.String("port", "8000", "Http port")
	addr       = flag.String("addr", "127.0.0.1", "IP address")
	prettyJSON = flag.Bool("pretty-json", false, "Pretty json in the stats")
)

// ImageServer implements main logic for serving images
type ImageServer struct {
	NumImages int
	WidthSum  int
	HeightSum int

	sync.RWMutex
}

// AddStats saves stat about image
func (is *ImageServer) AddStats(width int, height int) {

	is.Lock()
	defer is.Unlock()

	is.NumImages++
	is.WidthSum += width
	is.HeightSum += height
}

// ClearStats clear stats
func (is *ImageServer) ClearStats() {

	is.Lock()
	defer is.Unlock()

	is.NumImages = 0
	is.WidthSum = 0
	is.HeightSum = 0
}

// GetStats returns stats info
func (is *ImageServer) GetStats() (count int, avgWidth int, avgHeight int) {
	is.RLock()
	defer is.RUnlock()

	count = is.NumImages
	if is.NumImages > 0 {
		avgWidth = is.WidthSum / is.NumImages
		avgHeight = is.HeightSum / is.NumImages
	}

	return
}

// GetStatsJSON returns stats data in JSON format
func (is *ImageServer) GetStatsJSON() ([]byte, error) {
	var stat struct {
		NumImages   int `json:"num_images"`
		AvgWidthPx  int `json:"average_width_px"`
		AvgHeightPx int `json:"average_height_px"`
	}

	stat.NumImages, stat.AvgWidthPx, stat.AvgHeightPx = is.GetStats()

	var statsJSON []byte
	var err error
	if *prettyJSON {
		statsJSON, err = json.MarshalIndent(stat, "", "  ")
	} else {
		statsJSON, err = json.Marshal(stat)
	}
	return statsJSON, err
}

// GetStatsHTTP returns stat info for Http request
func (is *ImageServer) GetStatsHTTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	statsJSON, err := is.GetStatsJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsJSON)
}

// ClearStatsHTTP handles request for clear stats data
func (is *ImageServer) ClearStatsHTTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	is.ClearStats()
	is.GetStatsHTTP(w, r, nil)
}

// GetImageHTTP generate image of the certain type
func (is *ImageServer) GetImageHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	width, err := strconv.Atoi(ps.ByName("width_px"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	height, err := strconv.Atoi(ps.ByName("height_px"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blackColor := color.RGBA{0, 0, 0, 255}

	m := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			m.Set(x, y, blackColor)
		}
	}

	switch ps.ByName("img_type") {
	case "jpg":
		jpeg.Encode(w, m, nil)
	case "png":
		png.Encode(w, m)
	}
	is.AddStats(width, height)
}

func main() {

	flag.Parse()

	is := ImageServer{}
	router := httprouter.New()
	router.GET("/generate/:img_type/:width_px/:height_px", is.GetImageHTTP)
	router.GET("/stats", is.GetStatsHTTP)
	router.GET("/stats/clear", is.ClearStatsHTTP)

	connStr := fmt.Sprintf("%s:%s", *addr, *port)
	log.Fatal(http.ListenAndServe(connStr, router))
}
