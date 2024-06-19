package gameplay

import (
	"game/util"
	"log"
	"os"
)

const SETTINGS_FILEPATH string = "config"

type Settings struct {
	ScreenWidth, ScreenHeight int
}

func (g *Game) InitSettings()  {
	data, err := os.ReadFile(SETTINGS_FILEPATH)
	if err != nil {
		log.Println("[ReadFile]", err)
		g.ScreenWidth = 400
		g.ScreenHeight = 600
	}
	var settings Settings
	err = util.Decode(data, &settings)
	if err != nil {
		log.Println("[DecodeFile]", err)
		g.ScreenHeight = 400
		g.ScreenWidth = 600
	}
}

func (g *Game) SaveSettings() error {
	data, err := util.Encode(g.Settings)
	if err != nil {
		return err
	}
	file, err := os.Create(SETTINGS_FILEPATH)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	file.Close()
	return err
}
