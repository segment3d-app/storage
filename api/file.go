package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"archive/zip"

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
	isThumbnailGenerate := false

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

		uploadedFilePath := fmt.Sprintf("/files/%s",
			filepath.Join(filepath.Clean(folder), filepath.Base(file.Filename)))
		uploadedFiles = append(uploadedFiles, uploadedFilePath)

		// generate
		if !isThumbnailGenerate {
			if isVideo(file.Filename) {
				generateThumbnailForVideo(filePath)
				isThumbnailGenerate = true
			} else if isImage(file.Filename) {
				generateThumbnailForImage(filePath)
				isThumbnailGenerate = true
			}
		}
	}

	ctx.JSON(http.StatusOK, uploadFileResponse{Message: fmt.Sprintf("%d files uploaded successfully", len(files)), Url: uploadedFiles})
}

func generateThumbnailForImage(imagePath string) error {
	modifiedPath := strings.Replace(imagePath, "photos", "thumbnail", -1)
	sourceFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(strings.TrimRight(modifiedPath, filepath.Base(modifiedPath)), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories for thumbnail: %w", err)
	}

	destFile, err := os.Create(modifiedPath)
	if err != nil {
		return fmt.Errorf("failed to create destination image: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy image data: %w", err)
	}

	return nil
}

func generateThumbnailForVideo(videoPath string) error {
	modifiedPath := strings.Replace(videoPath, "photos", "thumbnail", -1)
	thumbnailPath := modifiedPath + ".jpg"
	if _, err := os.Stat(filepath.Dir(thumbnailPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(thumbnailPath), os.ModePerm)
	}
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01", "-frames:v", "1", thumbnailPath)
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}

func isVideo(filename string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mov", ".wmv"} // Add more as needed
	for _, ext := range videoExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

func isImage(filename string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"} // Add more as needed
	for _, ext := range imageExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
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

	filePath := filepath.Join(RootStorage, filepath.Clean("/"+capturedPath))

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if info.IsDir() {
		zipData, err := zipDirectory(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.Header("Content-Disposition", "attachment; filename="+info.Name()+".zip")
		ctx.Data(http.StatusOK, "application/zip", zipData)
		return
	}

	ctx.File(filePath)
}

func zipDirectory(srcDir string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type getThumbnailResponse struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

// @Summary Get file
// @Description Retrieve thumbnail from specified resource path
// @Tags file
// @Accept json
// @Produce json
// @Param path path string true "Path including any folders and subfolders to the file"
// @Success 200 {file} test "File retrieved successfully"
// @Router /thumbnail/{path} [get]
func (server *Server) getThumbnail(ctx *gin.Context) {
	capturedPath := ctx.Param("path")

	filePath := filepath.Join(RootStorage, filepath.Clean("/"+capturedPath))
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var firstFile string
	if info.IsDir() {
		firstFile, err = getFirstFileInDir(strings.Replace(filePath, "photos", "thumbnail", -1))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		folderPath := filepath.Dir(filePath)
		firstFile, err = getFirstFileInDir(strings.Replace(folderPath, "photos", "thumbnail", -1))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	url := fmt.Sprintf("/%s", firstFile)

	ctx.JSON(http.StatusAccepted, getThumbnailResponse{Url: url, Message: "thumnail image is successfully retrived"})
}

func getFirstFileInDir(dirPath string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {

			return filepath.Join(dirPath, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no files found in the directory")
}
