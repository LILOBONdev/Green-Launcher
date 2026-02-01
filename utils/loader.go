package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// URL MUST BE "GET"
func LoadFile(urlStr string, targetDir string, fileName string) error {
	if fileName == "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("неверный URL: %w", err)
		}
		fileName = path.Base(u.Path)
	}

	fullPath := filepath.Join(targetDir, fileName)

	if _, err := os.Stat(fullPath); err == nil {
		return nil
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер ответил %d", resp.StatusCode)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания директорий: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи: %w", err)
	}

	fmt.Printf("Загружено: %s\n", fileName)
	return nil
}
