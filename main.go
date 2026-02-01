package main

import (
	"minecraft_launcher/bin"
	"os/exec"

	"github.com/abemedia/go-webview"
	_ "github.com/abemedia/go-webview/embedded"
)

func main() {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("Microsoft Webview with Golang Example")
	w.SetSize(800, 600, webview.HintNone)

	w.SetTitle("Green Launcher")
	w.Navigate("C:/Users/NL/Desktop/Softs/minecraft_launcher_test/.frontend/main_page.html")

	err := w.Bind("play_minecraft", func(minecraft_version string, player_name string) {
		err_ := bin.Collect_Minecraft("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json", minecraft_version, player_name)
		if err_ != nil {
			panic(err_)
		}
	})

	if err != nil {
		panic(err)
	}

	err_ := w.Bind("start_minecraft", func(bat_path string) {
		cmd := exec.Command("cmd", "/C", bat_path)

		err := cmd.Run()

		if err != nil {
		}
	})

	if err_ != nil {
		panic(err_)
	}

	w.Run()
}
