package nightwatcher

const (
	nightWatcherLogTimeFormat = "2016-01-02 15:04:05"
)

type Reporter interface {
	Output(*Aim) error
}

var Reporters []Reporter = []Reporter{
	DefaultReporter,
}

func AddReporter(reporter Reporter) {
	Reporters = append(Reporters, reporter)
}

func SetReporter(reporter Reporter) {
	Reporters = []Reporter{
		reporter,
	}
}

func SetReporters(reporters []Reporter) {
	Reporters = reporters
}

func ClearReporters() {
	Reporters = []Reporter{}
}
