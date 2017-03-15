package stringModule

import (
	"time"

  "github.com/Sirupsen/logrus"
)

func loggingMiddleware(logger logrus.Entry) ServiceMiddleware {
	return func(next StringService) StringService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger logrus.Entry
	StringService
}

func (mw logmw) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.StringService.Uppercase(s)
	return
}

func (mw logmw) Count(s string) (output int, err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"method", "count",
			"input", s,
			"n", output,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.StringService.Count(s)
	return
}
