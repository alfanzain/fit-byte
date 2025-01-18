package models

type UploadImage struct {
	Uri string `json:"uri"`
}

type FileRequest struct {
	File string `json:"file" binding:"required"`
}
