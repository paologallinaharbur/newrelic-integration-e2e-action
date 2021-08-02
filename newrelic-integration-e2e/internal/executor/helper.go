package executor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func retry(log *logrus.Logger, attempts int, sleep time.Duration, f func() []error) error {

	var errors []error
	for i := 0; i < attempts; i++ {

		errors = f()
		if len(errors) == 0 {
			return nil
		}

		log.WithField("iteration", i).Warn("Error detected")
		for _, err := range errors {
			log.Error(err)
		}
		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, last errors: %v", attempts, errors)
}
