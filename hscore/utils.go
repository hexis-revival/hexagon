package hscore

import "net/http"

func GetMultipartFormValue(request *http.Request, key string) string {
	values, ok := request.MultipartForm.Value[key]
	if !ok || len(values) == 0 {
		return ""
	}

	return values[0]
}

func GetMultipartFormFile(request *http.Request, key string) []byte {
	files, ok := request.MultipartForm.File[key]
	if !ok || len(files) == 0 {
		return nil
	}

	file, err := files[0].Open()
	if err != nil {
		return nil
	}

	fileData := make([]byte, files[0].Size)
	_, err = file.Read(fileData)
	if err != nil {
		return nil
	}

	return fileData
}
