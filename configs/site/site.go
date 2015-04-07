package site

import (
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

// Individual Site Configuration Type
type SiteConfig struct {
	MongoDb   map[string]string
	Directory string
	Aliases   []string
	DbSession *mgo.Session
	Session	  string
}

// Constructor for SiteConfig
func New(file string) *SiteConfig {
	s := SiteConfig{
		MongoDb: make(map[string]string),
	}
	s.getSiteConfig(file)
	s.getDbConnection(file)
	return &s
}

// Get one site configuration and marshall it
func (s *SiteConfig) getSiteConfig(file string) {
	v := viper.New()
	v.SetConfigName(file)
	v.AddConfigPath(SITES_DIRECTORY)
	v.ReadInConfig()
	v.Marshal(s)
}

func (s *SiteConfig) getDbConnection(file) {
	u := "mongodb://" + s.MongoDb.user + ":" + s.MongoDb.password + "@" + s.MongoDb.host + "/" + s.MongoDb.db
	s.Db, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongodb server for %v, go error %v\n", file, err)
		os.Exit(1)
	}
}
