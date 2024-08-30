package converter

import (
	"fmt"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
)

func RunConversion(cfg *config.Configuration) error {
	var converter Converter

	switch cfg.ConversionType {
	case config.SrcConvert:
		converter = &SourceFileConverter{}
	case config.CfConvert:
		converter = &CfConverter{}
	default:
		return fmt.Errorf("неизвестный тип конвертации: %s", cfg.ConversionType)
	}

	if err := converter.Convert(cfg); err != nil {
		return fmt.Errorf("ошибка конвертации: %w", err)
	}
	return nil
}
