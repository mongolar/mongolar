package site

import (
	"github.com/jasonrichardsmith/mongolar/logger"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

// Individual Site Configuration Type
type SiteConfig struct {
	MongoDb           map[string]string      //Configuration for MongoDB Connection
	Directory         string                 // Directory for html and assets
	Aliases           []string               // Site Aliases/Domains
	SessionExpiration int64                  // When to expire a users Session
	TemplateEndpoint  string                 // URL where will be stored
	ForeignDomains    []string               // This will whitelist domains for loading assets from other domains
	AngularModules    []string               // A slice of angularjs modules to load
	PublicValues      map[string]string      // These values can be directly invoked from the domain controller
	Misc              map[string]interface{} // Where you can store any other value not defined here
	Logger            logger.LogChannel      // A channeel for writing Logs
	DbSession         *mgo.Session           // The master MongoDb session that gets copied
	FourOFour         string
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

// Establish a Database connection and attach it to the site configuration
func (s *SiteConfig) getDbConnection(f string) {
	u := "mongodb://" + s.MongoDb["user"] + ":" + s.MongoDb["password"] + "@" + s.MongoDb["host"] + "/" + s.MongoDb["db"]
	var err error
	s.DbSession, err = mgo.Dial(u)
	if err != nil {
		log.Fatal(err)
	}
}

// Attach a logger channel to log errors predictably.
func (s *SiteConfig) getLogger() {
	s.Logger = logger.New(s.DbSession)
}
