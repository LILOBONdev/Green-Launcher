package utils

const (
	// -Djava.library.path="minecraft\versions\1.12.2\natives"
	BASH_ARGS = `
@echo off
chcp 65001 >nul

cd /d "%~dp0.."
title Minecraft {mc_ver}
setlocal enabledelayedexpansion

set "REL_PATH=.\minecraft\jre\x64\bin\java.exe"
for %%I in ("%REL_PATH%") do set "J=%%~fI"

set GAME_JAR=minecraft\minecraft\versions\{mc_ver}\{mc_ver}.jar

echo [INFO] Starting Minecraft {mc_ver}

echo -cp > java_args.txt
<nul set /p ="%GAME_JAR%" >> java_args.txt

for /r "minecraft\minecraft\libraries" %%i in (*.jar) do (
	<nul set /p =";%%i" >> java_args.txt
)

"%J%" -Xmx2G @java_args.txt net.minecraft.client.main.Main ^
	--username {plr_name} ^
	--version {mc_ver} ^
	--gameDir minecraft\minecraft ^
	--assetsDir minecraft\minecraft\assets ^
	--assetIndex {asst_idx} ^
	--accessToken 0 ^
	--userType msa

if %errorlevel% neq 0 (
	echo.
	echo [ERR] Exit Code: %errorlevel%
)

RD /S /Q "minecraft/minecraft/libraries"
pause
	`

	Blue_color  = "\033[34m"
	Green_color = "\033[32m"
	Reset_color = "\033[0m"
)
