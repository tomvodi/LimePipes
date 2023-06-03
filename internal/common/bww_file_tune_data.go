package common

import "github.com/rs/zerolog/log"

type tuneFileData struct {
	Title string
	Data  []byte
}

type BwwFileTuneData struct {
	tuneData []tuneFileData
}

func (b *BwwFileTuneData) TuneTitles() (titles []string) {
	for _, tuneData := range b.tuneData {
		titles = append(titles, tuneData.Title)
	}
	return titles
}

func (b *BwwFileTuneData) Data(idx int) []byte {
	if idx >= len(b.tuneData) {
		log.Error().Msgf("no tune data for index found")
		return nil
	}

	return b.tuneData[idx].Data
}

func (b *BwwFileTuneData) AddTuneData(title string, data []byte) {
	tuneData := tuneFileData{
		Title: title,
		Data:  data,
	}
	b.tuneData = append(b.tuneData, tuneData)
}

func NewBwwFileTuneData() *BwwFileTuneData {
	return &BwwFileTuneData{
		tuneData: make([]tuneFileData, 0),
	}
}
