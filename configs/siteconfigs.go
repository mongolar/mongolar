package configs

import (
	"github.com/Sirupsen/logrus"
	"github.com/mongolar/mongolar/logger"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"sort"
	"time"
)

// Individual Site Configuration Type
type SiteConfig struct {
	MongoDb            map[string]string //Configuration for MongoDB Connection
	Directory          string            // Directory for html and assets
	Aliases            []string          // Site Aliases/Domains
	SessionExpiration  time.Duration     // When to expire a users Session
	TemplateEndpoint   string            // URL where will be stored
	ForeignDomains     []string          // This will whitelist domains for loading assets from other domains
	AngularModules     []string          // A slice of angularjs modules to load
	PublicValues       map[string]string // These values can be directly invoked from the domain controller
	FourOFour          string
	APIEndPoint        string
	Controllers        []string
	ElementControllers []string
	Logger             *logrus.Logger // Logrus logger
	DbSession          *mgo.Session   // The master MongoDb session that gets copied
	RawConfig          *viper.Viper
}

// Constructor for SiteConfig
func NewSiteConfig(f string) *SiteConfig {
	s := SiteConfig{
		MongoDb: make(map[string]string),
	}
	s.getSiteConfig(f)
	s.getDbConnection(f)
	s.getLogger(f)
	sort.Strings(s.Controllers)
	return &s
}

// Get one site configuration and marshall it
func (s *SiteConfig) getSiteConfig(file string) {
	v := viper.New()
	v.SetConfigName(file)
	v.AddConfigPath(ServerConfig.SitesDirectory)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = v.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	s.RawConfig = v
}

// Establish a Database connection and attach it to the site configuration
func (s *SiteConfig) getDbConnection(f string) {
	u := "mongodb://" + s.MongoDb["user"] + ":" + s.MongoDb["password"] + "@" + s.MongoDb["host"] + "/" + s.MongoDb["db"]
	dbs, err := mgo.Dial(u)
	if err != nil {
		log.Fatal(err)
	}
	s.DbSession = dbs
}

// Attach a logger channel to log errors predictably.
func (s *SiteConfig) getLogger(f string) {
	s.Logger = logger.New(ServerConfig.LogDirectory + f)
}
