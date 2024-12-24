package converter

import "github.com/firstBitSportivnaya/files-converter/internal/config"

type Converter interface {
	Convert(cfg *config.Configuration) error
}
