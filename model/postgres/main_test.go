package postgres

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/DHowett/ghostbin/model"
	"github.com/Sirupsen/logrus"
)

type noopChallengeProvider struct{}

func (n *noopChallengeProvider) DeriveKey(string, []byte) []byte {
	return []byte{'a'}
}

func (n *noopChallengeProvider) RandomSalt() []byte {
	return []byte{'b'}
}

func (n *noopChallengeProvider) Challenge(message []byte, key []byte) []byte {
	return append(message, key...)
}

var gTestProvider model.Provider

func TestMain(m *testing.M) {
	flag.Parse()
	sqlDb, err := sql.Open("postgres", "postgresql://ghostbin:password@localhost/ghostbintest?sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = sqlDb.Exec(`
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO ghostbin;
	GRANT ALL ON SCHEMA public TO public;
	`)
	fmt.Println(err)

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	if testing.Verbose() {
		logger.Level = logrus.DebugLevel
	}
	gTestProvider, err = model.Open("postgres", sqlDb, &noopChallengeProvider{}, model.FieldLoggingOption(logger))
	if err != nil {
		panic(err)
	}
	e := m.Run()
	os.Exit(e)
}
