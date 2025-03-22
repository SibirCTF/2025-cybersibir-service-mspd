package utils

import (
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"sibir2025/service/app_logic/models"

	"github.com/nfnt/resize"
)

func Upload_sus_img(w http.ResponseWriter, r *http.Request, sus models.Suspect) {
	uploaded_img := r.MultipartForm.File["sus_img"]
	if len(uploaded_img) == 0 {
		return
	}
	fileHeader := uploaded_img[0]
	file, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}
	resizedImg := resize.Resize(200, 200, img, resize.Lanczos3)
	dir := "./static/sus"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	filename := fmt.Sprintf("%d.png", sus.ID)
	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return
	}
	defer dst.Close()
	if err := png.Encode(dst, resizedImg); err != nil {
		return
	}
}
