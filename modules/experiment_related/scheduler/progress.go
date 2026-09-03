package scheduler

import "github.com/schollz/progressbar/v3"

func newProgressBar(max int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		max,
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowBytes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(100),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)
}
