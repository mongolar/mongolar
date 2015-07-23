// Site configurations load configuration files for each site.
// It has some required set paramaters it loads from the Yaml file, it also
// loads tools with site specific settings, specifically:
// Logger - Is set with the parameters needed to log to the sites specific logs
// 	Logger is an instance of Logrus
// DbSession - Is the original Db connection for the site.  The request wrapper
// 	will make a copy locally for all requests and close those connections.
//	Read more under Wrapper.
// RawConfig - is an Instance of Viper for the configuration, so the user can
//	retrieve values not defeined in the core SiteConfig structure.  This can
//	be done by calling SiteConfig.RawConfig.Get('my_value') or
//	SiteConfig.RawConfig.MarshalKey('my_value', my_type)

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
// 	MongoDb: Configuration for MongoDB Connection
// 	Directory: Directory for html and assets
// 	AssetsDirectory: Directory for assets, default is assets
// 	Aliases: Site Aliases/Domains
// 	SessionExpiration: When to expire a users Session
// 	TemplateEndpoint: URL where will be stored
// 	ForeignDomains: This will whitelist domains for loading assets from other
//		domains
// 	AngularModules: A slice of angularjs modules to load
// 	PublicValues: These values can be directly invoked from the domain controller
// 	FourOFour: Required for redirects to all pages not found.
// 	APIEndPoint: Where all api requests are destined. Any valid url string(not /)
// 	Controllers: This is a list of valid controller map end points that are
//		valid.  This is so you can establish per site functionality
// 	ElementControllers: Elements availabled to be created in the UI
// 	Logger:	Logrus logger
// 	DbSession: The master MongoDb session that gets copied
// 	RawConfig: Raw viper configuration

type SiteConfig struct {
	MongoDb            map[string]string
	Directory          string
	AssetsDirectory	   string
	Aliases            []string
	SessionExpiration  time.Duration
	TemplateEndpoint   string
	ForeignDomains     []string
	AngularModules     []string
	PublicValues       map[string]string
	FourOFour          string
	APIEndPoint        string
	Controllers        []string
	ElementControllers []string
	Logger             *logrus.Logger
	DbSession          *mgo.Session
	RawConfig          *viper.Viper
}

// Constructor for SiteConfig, takes config filename as an argument.
func NewSiteConfig(f string) *SiteConfig {
	s := SiteConfig{
		MongoDb: make(map[string]string),
		AssetsDirectory: "assets"
	}
	// Marshall config based on filename
	s.getSiteConfig(f)
	s.getDbConnection()
	// Set log file based on config filename
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
func (s *SiteConfig) getDbConnection() {
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
