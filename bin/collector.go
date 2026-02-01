package bin

import (
	"fmt"
	"io"
	"log"
	"minecraft_launcher/utils"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func MakeBat(minecraftVersion, playerName, assetIndex string) {
	outputDir := "Launch"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Не удалось создать директорию: %v", err)
		return
	}

	assetIndex = strings.Replace(assetIndex, ".json", "", -1)
	replacer := strings.NewReplacer(
		"{mc_ver}", minecraftVersion,
		"{plr_name}", playerName,
		"{asst_idx}", assetIndex,
	)
	batCode := replacer.Replace(utils.BASH_ARGS)
	encoder := charmap.Windows1251.NewEncoder()
	win1251Code, err := encoder.String(batCode)
	if err != nil {
		log.Printf("Ошибка смены кодировки: %v", err)
		win1251Code = batCode
	}

	fileName := fmt.Sprintf("%s/Minecraft %s [%s].bat", outputDir, minecraftVersion, playerName)
	err = os.WriteFile(fileName, []byte(win1251Code), 0644)
	if err != nil {
		log.Printf("Ошибка записи .bat: %v", err)
		return
	}

	fmt.Printf("Клиент %s успешно создан!\n", fileName)
}

func Collect_Minecraft(minecraft_manifest_url string, minecraft_version string, player_name string) error {
	minecraft_versions := make(map[string]string)
	resp, err := http.Get(minecraft_manifest_url)

	if err != nil {
		return fmt.Errorf("Ошибка при загрузке основного манифеста Minecraft: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	encoded_body := utils.JsonFormaters(body)

	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела манифеста Minecraft: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка в ходе загрузки основного манифеста Minecraft: %d", resp.StatusCode)
	}

	if versions, ok := encoded_body["versions"].([]interface{}); ok {
		for _, version_block := range versions {
			if version_body, ok := version_block.(map[string]interface{}); ok {
				minecraft_versions[version_body["id"].(string)] = version_body["url"].(string)
			}
		}
	}

	mcversion, ok := minecraft_versions[minecraft_version]
	if ok {
	} else {
		return fmt.Errorf("Версии %s нет в манифесте Minecraft", mcversion)
	}

	api_url := minecraft_versions[minecraft_version]

	Load_client(api_url)
	Load_libraries(api_url, minecraft_version)

	err, asset_index_url := Load_resources(api_url)
	if err != nil {
		return fmt.Errorf("Ошибка при загрузке ресурсов Minecraft: %w", err)
	}
	splitted_url := strings.Split(asset_index_url, "/")
	asset_indexx := splitted_url[len(splitted_url)-1]

	MakeBat(minecraft_version, player_name, asset_indexx)

	return nil
}
