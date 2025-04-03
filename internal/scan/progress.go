package scan

import "github.com/jedib0t/go-pretty/v6/progress"

// initializeProgress sets up the progress writer and tracker.
func initializeProgress(totalIPs int) (*progress.Writer, *progress.Tracker) {
	pw := progress.NewWriter()
	pw.SetAutoStop(true)
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowValue(true)
	pw.SetStyle(progress.StyleBlocks)
	go pw.Render()

	tracker := &progress.Tracker{
		Message: "Scanning...",
		Total:   int64(totalIPs),
		Units:   progress.UnitsDefault,
	}
	pw.AppendTracker(tracker)

	return &pw, tracker
}
