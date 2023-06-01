package common

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

func (b *BwwFileTuneData) HasDataForTune(title string) bool {
	for _, tuneData := range b.tuneData {
		if tuneData.Title == title {
			return true
		}
	}
	return false
}

func (b *BwwFileTuneData) DataForTune(title string) []byte {
	for _, tuneData := range b.tuneData {
		if tuneData.Title == title {
			return tuneData.Data
		}
	}
	return nil
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
