package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const MAX_FILE_SIZE = 10 << 20 // 10MB
const DEFAULT_ERROR_CODE = 1
const STORAGE_SONG_PATH = "/song"

type UploadSongResponse struct {
	ResourceID   string `json:"resource_id"`
	ResourceLink string `json:"resource_link"`
}

func (h *Handler) UploadSong(w http.ResponseWriter, r *http.Request, params map[string]string) {
	r.ParseMultipartForm(MAX_FILE_SIZE)
	file, handler, err := r.FormFile("file")

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, DEFAULT_ERROR_CODE, fmt.Sprintf("could not retrieve file: %v", err))
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, file); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, DEFAULT_ERROR_CODE, fmt.Sprintf("could not read file content: %v", err))
		return
	}

	fileName, fileLink, err := h.storageService.CreateFile(buf.Bytes(), handler.Filename, STORAGE_SONG_PATH)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, DEFAULT_ERROR_CODE, fmt.Sprintf("could not create file: %v", err))
		return
	}

	sendOkResponse(w, UploadSongResponse{ResourceID: fileName, ResourceLink: fileLink})
}
