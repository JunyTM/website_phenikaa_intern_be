package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

const CHAR_SET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateFilename(prefix string, ext string, length int) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = CHAR_SET[seededRand.Intn(len(CHAR_SET))]
	}
	return prefix + "_" + string(b) + "." + ext
}

func UploadAndReturnFilePath(rootPath string, prefixDir string, storageDir string, prefixPath string, userID int, file multipart.File, headers *multipart.FileHeader) (string, error) {
	// log.Println(headers.Filename)
	fileNameArr := strings.Split(headers.Filename, ".")
	if len(fileNameArr) < 2 {
		// log.Println(fileNameArr)
		return "fail", errors.New("file name is not acceptable " + headers.Filename)
	}
	extension := fileNameArr[len(fileNameArr)-1]

	prefixPathList := strings.Split(prefixPath, ".")
	if len(prefixPathList) > 1 {
		prefixPath = prefixPathList[0]
	}

	rootAndStoragePath := rootPath + "/" + storageDir
	if _, err := os.Stat(rootAndStoragePath + "/" + prefixDir); os.IsNotExist(err) {
		if err := os.Mkdir(rootAndStoragePath+"/"+prefixDir, 0755); err != nil {
			return "", err
		}
	}

	fileDir := rootAndStoragePath + "/" + prefixDir + "/" + strconv.Itoa(userID)
	fileName, err := UploadFile(prefixPath, fileDir, file, extension)
	if err != nil {
		log.Println(err)
		return "", err
	}

	filePath := storageDir + "/" + prefixDir + "/" + strconv.Itoa(userID) + "/" + fileName
	return filePath, nil
}

func UploadFile(prefix string, dirPath string, file multipart.File, extension string) (string, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return "", errors.New("do not permission create upload folder")
		}
	}

	fileName := GenerateFilename(prefix, extension, 16)
	filePath := dirPath + "/" + fileName
	destFile, err := os.Create(filePath)
	if err != nil {
		log.Println("true err:", err)
		return "", errors.New("do not permission to create new file")
	}
	defer destFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errors.New("can't read file content")
	}

	if _, err := destFile.Write(fileBytes); err != nil {
		return "", errors.New("Do not permission to  upload file")
	}

	return fileName, nil
}
