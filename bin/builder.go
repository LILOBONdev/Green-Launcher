package bin

import (
	"fmt"
	"io"
	"minecraft_launcher/utils"
	"net/http"
	"strings"
	"sync"
)

func Load_client(minecraft_api_url string) error {
	resp, err := http.Get(minecraft_api_url)
	api_splitted := strings.Split(minecraft_api_url, "/")
	minecraft_version := api_splitted[len(api_splitted)-1]
	minecraft_version = strings.Replace(minecraft_version, ".json", "", -1)

	if err != nil {
		return fmt.Errorf("Ошибка при загрузке клиента Minecraft: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	encoded_body := utils.JsonFormaters(body)

	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела манифеста Minecraft: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка в ходе выполнения загрузки клиента Minecraft: %d", resp.StatusCode)
	}

	utils.LoadFile(minecraft_api_url, fmt.Sprintf("minecraft/minecraft/versions/%s/", minecraft_version), "")

	if downloads, ok := encoded_body["downloads"].(map[string]interface{}); ok {
		if client, ok := downloads["client"].(map[string]interface{}); ok {
			if url, ok := client["url"].(string); ok {
				current_path := fmt.Sprintf("minecraft/minecraft/versions/%s/", minecraft_version)
				utils.LoadFile(url, current_path, fmt.Sprintf("%s.jar", minecraft_version))
			}
		}
	}
	return nil
}

func Load_libraries(minecraft_api_url string, minecraft_version string) error {
	resp, err := http.Get(minecraft_api_url)
	if err != nil {
		return fmt.Errorf("Ошибка при загрузке библиотек Minecraft: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	encoded_body := utils.JsonFormaters(body)

	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела манифеста Minecraft: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка в ходе выполнения загрузки библиотек: %d", resp.StatusCode)
	}

	if libs, ok := encoded_body["libraries"].([]interface{}); ok {
		for _, item := range libs {
			if lib, ok := item.(map[string]interface{}); ok {
				if downloads, found := lib["downloads"].(map[string]interface{}); found {
					if artifact, found := downloads["artifact"].(map[string]interface{}); found {
						splitted_path := strings.Split(artifact["path"].(string), "/")

						ppath := strings.Join(splitted_path[:len(splitted_path)-1], "/")
						current_path := "minecraft/minecraft/libraries/" + ppath
						utils.LoadFile(artifact["url"].(string), current_path, "")
					}

					if classifiers, found := downloads["classifiers"].(map[string]interface{}); found {
						if natives, found := classifiers["natives-windows"].(map[string]interface{}); found {
							if natives_url, found := natives["url"].(string); found {
								utils.UnzipJarFromURL(natives_url, fmt.Sprintf("minecraft/minecraft/versions/%s/natives", minecraft_version))
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func Load_resources(minecraft_api_url string) (error, string) {
	resp, err := http.Get(minecraft_api_url)
	if err != nil {
		return fmt.Errorf("Ошибка при загрузке манифеста: %w", err), ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка манифеста, статус: %d", resp.StatusCode), ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела манифеста: %w", err), ""
	}

	encoded_body := utils.JsonFormaters(body)

	assetIndex, ok := encoded_body["assetIndex"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("assetIndex не найден"), ""
	}

	resources_url, ok := assetIndex["url"].(string)
	if !ok {
		return fmt.Errorf("url ресурсов не найден"), ""
	}

	resources_resp, err := http.Get(resources_url)
	if err != nil {
		return fmt.Errorf("Ошибка при загрузке индекса ресурсов: %w", err), ""
	}
	defer resources_resp.Body.Close()

	resources_body, err := io.ReadAll(resources_resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела ресурсов: %w", err), ""
	}

	utils.LoadFile(resources_url, "minecraft/minecraft/assets/indexes/", "")

	assets_data := utils.JsonFormaters(resources_body)

	if objects, ok := assets_data["objects"].(map[string]interface{}); ok {
		fmt.Printf("Найдено ресурсов: %d. Начинаю загрузку...\n", len(objects))

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 20)

		for name, data := range objects {
			resourceInfo, ok := data.(map[string]interface{})
			if !ok {
				continue
			}

			hash, ok := resourceInfo["hash"].(string)
			if !ok {
				continue
			}

			prefix := hash[0:2]
			downloadUrl := fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", prefix, hash)
			destPath := fmt.Sprintf("minecraft/minecraft/assets/objects/%s/", prefix)

			wg.Add(1)
			go func(url string, path string, fileName string) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				err := utils.LoadFile(url, path, hash)
				if err != nil {
					fmt.Printf("Ошибка загрузки %s: %v\n", fileName, err)
				}
			}(downloadUrl, destPath, name)
		}

		wg.Wait()
	}

	return nil, resources_url
}
