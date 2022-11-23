package grpc

import (
	"fmt"
	"net/http"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
)

type UploadSongResponse struct {
	ResourceID   string `json:"resource_id"`
	ResourceLink string `json:"resource_link"`
}

func (h *Handler) UploadSong(w http.ResponseWriter, r *http.Request, params map[string]string) {
	r.ParseMultipartForm(constants.MAX_FILE_SIZE)
	file, header, err := r.FormFile("file")

	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, constants.DEFAULT_ERROR_CODE, fmt.Sprintf("no file is found: %v", err))
		return
	}
	defer file.Close()

	resourceID, resourceLink, err := h.songService.UploadSong(r.Context(), header, file)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, constants.DEFAULT_ERROR_CODE, err.Error())
		return
	}

	sendOkResponse(w, UploadSongResponse{ResourceID: resourceID, ResourceLink: resourceLink})
}