package controller

import (
	"fmt"
	"net/http"
	"phenikaa/infrastructure"
	"phenikaa/model"
	"phenikaa/utils"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type FileController interface {
	UploadFileWithPath(w http.ResponseWriter, r *http.Request)
}

type fileController struct {
}

// Upload file godoc
// @tags        file-manager-apis
// @Summary     Update file
// @Description Update file with path
// @Accept      json
// @Produce     json
// @Param       studentId   	   	query    integer true "student id"
// @Param       prefixPath query    string  true "prefix path"
// @Param       prefixDir  query    string  true "prefix dir"
// @Param       file       formData file    true "multi image file"
// @Security    ApiKeyAuth
// @Success     200 {object} controller.Response
// @Router      /file/upload-with-path [post]
func (f *fileController) UploadFileWithPath(w http.ResponseWriter, r *http.Request) {
	var res Response
	queryValues := r.URL.Query()
	studentId, err := strconv.Atoi(queryValues.Get("studentId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		res = Response{
			Data:    nil,
			Message: "Submit student failed: " + err.Error(),
			Success: false,
		}
		render.JSON(w, r, res)
		return
	}
	prefixPath := queryValues.Get("prefixPath")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		res = Response{
			Data:    nil,
			Message: "Submit student failed: " + err.Error(),
			Success: false,
		}
		render.JSON(w, r, res)
		return
	}
	prefixDir := queryValues.Get("prefixDir")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, http.StatusText(400), 400)
		res = Response{
			Data:    nil,
			Message: "Submit student failed: " + err.Error(),
			Success: false,
		}
		render.JSON(w, r, res)
		return
	}

	storageDirPath := infrastructure.GetStoragePath()
	rootPath := infrastructure.GetRootPath()

	//25mb
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		return
	}

	// hoc ba
	formData := r.MultipartForm
	listFiles := formData.File["file"]

	// Cập nhập lại đường dẫn file vào trong database
	db := infrastructure.GetDB()
	var fileUrls []string
	// var media []model.Recruitment
	for i := range listFiles {
		file, err := listFiles[i].Open()

		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer file.Close()
		if listFiles[i].Filename != "" {
			if prefixPath != "" {
				prefixPath = prefixPath + "_" + listFiles[i].Filename
			} else {
				prefixPath = listFiles[i].Filename
			}
		}
		fileUrl, err := utils.UploadAndReturnFilePath(rootPath, prefixDir, storageDirPath, prefixPath, studentId, file, listFiles[i])
		if fileUrl == "fail" {
			res = Response{
				Data:    nil,
				Message: "Upload failed; msg=" + err.Error(),
				Success: false,
			}
			render.JSON(w, r, res)
			return
		}
		fileUrls = append(fileUrls, fileUrl)
		// fileUrls = strings.Split(fileUrl, ",")
		// for _, fileUrl := range fileUrls {
		// 	if fileUrl == "" {
		// 		continue
		// 	}
				
		// 	}
		// }


	}
	fileString := strings.Join(fileUrls, ",")
	if err := db.Model(&model.Recruitment{}).Where("profile_id = (?)", uint(studentId)).Update("profile_path", fileString).Error; err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
	

	res = Response{
		Data:    fileString,
		Message: "upload file successful",
		Success: true,
	}
	render.JSON(w, r, res)
}


func NewFileController() FileController {
	return &fileController{}
}