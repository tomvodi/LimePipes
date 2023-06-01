package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/interfaces"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
)

const matchStart = 0

type bwwFileSplitter struct {
}

func (b *bwwFileSplitter) SplitFileData(data []byte) (fileTuneData *common.BwwFileTuneData, err error) {
	fileTuneData = common.NewBwwFileTuneData()

	reg := regexp.MustCompile(`"[^"]*"\s*,\s*\(\s*T`)
	indexes := reg.FindAllIndex(data, -1)
	results := make([][]byte, len(indexes))
	for i, element := range indexes {
		if i+1 == len(indexes) {
			// last element
			results[i] = data[element[matchStart]:]
		} else {
			nextElement := indexes[i+1]
			results[i] = data[element[matchStart]:nextElement[matchStart]]
		}
	}

	re := regexp.MustCompile(`"[^"]*"`)
	for _, tune := range results {
		titles := re.FindSubmatch(tune)
		if len(titles) != 1 {
			log.Error().Msgf("tune has more than one title")
		}
		if len(titles) == 0 {
			msg := "no title found in tune"
			log.Error().Msgf(msg)
			return nil, fmt.Errorf(msg)
		}

		fileTuneData.AddTuneData(
			string(bytes.Trim(titles[0], `"`)),
			tune,
		)
	}

	return fileTuneData, nil
}

func NewBwwFileTuneSplitter() interfaces.BwwFileByTuneSplitter {
	return &bwwFileSplitter{}
}
