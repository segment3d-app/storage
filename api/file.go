package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var RootStorage = "./files"

type uploadFileResponse struct {
	Message string   `json:"message"`
	Url     []string `json:"url"`
}

// uploadFile uploads a file to a specified folder and returns the URL.
// @Summary Upload file
// @Description Uploads a file to the specified folder within the server's storage directory.
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param folder formData string true "Folder where the file will be uploaded"
// @Param file formData []file true "File(s) to upload"
// @Success 200 {object} uploadFileResponse "Upload file success"
// @Router /upload [post]
func (server *Server) uploadFile(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1000<<20) // Limit the request body to 1GB

	folder := ctx.PostForm("folder")
	form, _ := ctx.MultipartForm()
	files := form.File["file"]

	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no file provided")))
		return
	}

	var uploadedFiles []string

	for _, file := range files {
		filePath := filepath.Join(RootStorage, filepath.Clean(folder), filepath.Base(file.Filename))

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		uploadedFilePath := fmt.Sprintf("%s://%s:%s/files/%s",
			server.config.StorageProtocol,
			server.config.StorageAddress,
			server.config.StoragePort,
			filepath.Join(filepath.Clean(folder), filepath.Base(file.Filename)))
		uploadedFiles = append(uploadedFiles, uploadedFilePath)
	}

	ctx.JSON(http.StatusOK, uploadFileResponse{Message: fmt.Sprintf("%d files uploaded successfully", len(files)), Url: uploadedFiles})
}

// @Summary Get file
// @Description Retrieve file data from specified path within the server's storage directory.
// @Tags file
// @Accept json
// @Produce octet-stream
// @Param path path string true "Path including any folders and subfolders to the file"
// @Success 200 {file} file "File retrieved successfully"
// @Router /files/{path} [get]
func (server *Server) getFile(ctx *gin.Context) {
	capturedPath := ctx.Param("path")

	filePath := filepath.Join(RootStorage, filepath.Clean("/"+capturedPath)) // Prepending slash to ensure path is correctly joined

	if info, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if info.IsDir() {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.File(filePath)
}
