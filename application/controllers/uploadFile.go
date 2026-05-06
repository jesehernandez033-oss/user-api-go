package controllers

import (
	"awesomeProject1/infrastructure"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// UploadFile godoc
// @Summary Subir archivo
// @Description Sube un archivo al servidor
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Archivo a subir"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /uploadFile [post]
func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UploadFile | error=invalid multipart form | err=%v",
			err,
		)
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=UploadFile | error=missing or invalid file | err=%v",
			err,
		)
		http.Error(w, "error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UploadFile | error=creating uploads folder | err=%v",
			err,
		)
		http.Error(w, "error preparing upload folder", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UploadFile | filename=%s | error=saving file | err=%v",
			handler.Filename,
			err,
		)
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UploadFile | filename=%s | error=writing file | err=%v",
			handler.Filename,
			err,
		)
		http.Error(w, "error writing file", http.StatusInternalServerError)
		return
	}

	infrastructure.Logger.Printf(
		"INFO | file uploaded | filename=%s | size=%d",
		handler.Filename,
		handler.Size,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "file uploaded successfully",
	})
}
