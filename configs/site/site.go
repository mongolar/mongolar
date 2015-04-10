package site

import (
	"fmt"
	"github.com/jasonrichardsmith/mongolar/logger"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"os"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

// Individual Site Configuration Type
type SiteConfig struct {
	MongoDb           map[string]string
	Directory         string
	Aliases           []string
	SessionExpiration int64
	TemplateEndpoint  string
	ForeignDomains    []string
	AngularModules    []string
	PublicValues      map[string]string
	Logger            logger.LogChannel
	DbSession         *mgo.Session
}

// Constructor for SiteConfig
func New(f string) *SiteConfig {
	s := SiteConfig{
		MongoDb: make(map[string]string),
	}
	s.getSiteConfig(f)
	s.getDbConnection(f)
	s.getLogger()
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

func (s *SiteConfig) getDbConnection(f string) {
	u := "mongodb://" + s.MongoDb["user"] + ":" + s.MongoDb["password"] + "@" + s.MongoDb["host"] + "/" + s.MongoDb["db"]
	var err error
	s.DbSession, err = mgo.Dial(u)
	if err != nil {
		fmt.Printf("Can't connect to mongodb server for %v, go error %v\n", f, err)
		os.Exit(1)
	}
}

func (s *SiteConfig) getLogger() {
	s.Logger = logger.New(s.DbSession)
}
