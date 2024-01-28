package infrastructure

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

func WriteToFile(filename, text string) error {
	// Открываем файл для записи. Если файл не существует, он будет создан.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(text)
	if err != nil {
		return err
	}

	err = writer.Flush() // Важно вызвать Flush, чтобы убедиться, что все данные записаны в файл.
	if err != nil {
		return err
	}
	return nil
}

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// downloadFile загружает файл из URL и сохраняет его по указанному пути
func downloadFile(filepath, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
