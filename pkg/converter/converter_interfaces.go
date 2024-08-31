package converter

import "github.com/firstBitSportivnaya/files-converter/pkg/config"

type Converter interface {
	Convert(cfg *config.Configuration) error
}
