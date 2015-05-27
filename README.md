#Mongolar

Mongolar is a AngularJS driven scalable Content Management System written in AngularJS and Go with MongoDB as a backend.

##Basics
###AngularJS
Essentialy AngularJS bootstraps itself upon page load.  It then begins to process any "mongolar" directives, loading content to the page.

Mongolar directives in AngularJS perform api calls based on settings in the directive that give the address for the call, an id, a template and a dynamic id (an element id that makes dynamic page content loading easier).

The simplest implementation for a directive would be to load page elements from a path and then the returned values would then create more directives to subsequently make more API requests

###Go

The Mongolar binary will work as a webserver/API server.

Mongolar serves along several different paths
 - /assets : Holds all html, css, and js assets
 - /mongolarconfig.js : This is a site specific js file that loads all the mongolar-js specific configurations for a particular site
 - /apiendpoint : This a per site configurable endpoint for api requests.  This can be any url friendly string
 - default: If all the above routing does not work then the index.html will be served allowing AngularJS to take over routing.

###MongoDB
MongoDb stores stuff like NoSQL Dbs do.

##Scalable

This package seeks to achieve scalability in several ways.

Each Mongolar server instance is stateless.  You can spin up as many servers as you wish behind a load balancer, and there is not extra configuration.

All html is cacheable and can be served through a cdn.

No templating is performed on the server.

Every web request is a microtransaction for individual pieces of content, vs one monoloithic request for a single web page.  So intensive processes do not hold up the entire page load.

Every piece of content (read API request) is individually addressable so if situated behind a tool like varnish, you can easily cache content without making requests to the system.

MongoDB seems to scale rather well, at least for these purposes.

## Demo
I plan to release a demo for this server shortly.  In the demo you will be able to manipulate content and view the results

##Setup

To setup the go server you will do several things (assuming go is already setup on your system):

###1. Download mongolar.
You can clone this repository anywhere you wish.

###2. Configuration
Mongolar currently only has one configuration file for a server and as many site files as you wish

All configurations are in YAML format

The default location for a server config would be here "/etc/mongolar", if this does not work for you you can set the "MONGOLAR_SERVER_CONFIG" to the folder you want to use.

Example:
```yaml
"Port": "80"
"SitesDirectory": "/etc/mongolar/enabled"
"LogDirectory": "/var/log/mongolar/"
```

Based on the server config Mongolar will attempt to load all site configuration files into memory, so given the above configuration you would create a yaml:

"/etc/mongolar/enabled/my_site.yaml" <--yaml suffix required

Ideally you would create your configs in /etc/mongolar/available/my_site.yaml and then symlink to the file, making configs easier to manage.

Example:
```yaml
MongoDb:
        user: my_db_user
        password: my_password
        host: "db_domain:12345"
        db: my_db
Directory: "/my/files/directory"
Aliases:
        - mongolar.org
        - test-1.com
        - test-1.1.com
PublicValues:
        "test" : "Test Value"
FourOFour: "page_not_found"
AngularModules:
        - "formly"
        - "formlyBootstrap"
        - "angular-growl"
        - "ui.bootstrap"
        - "ngSanitize"
        - "ui.sortable"
TemplateEndpoint: "assets/templates"
SessionExpiration: 10
APIEndPoint: "my_end_point"
Controllers:
        - "path"
        - "content"
        - "wrapper"
        - "slug"
        - "admin"
        - "domain_public_value"
        - "login"
OAuthLogins:
        "github":
                "client_id": "client_id_here"
                "client_secret": "client_secret_here"
                "login_text": "Login with github"
LoginURLs: 
        "login": "/login"
        "success" : "/login/success"
        "failure": "/login/failure"
        "access_denied": "/access_denied"
```
###3. Frontend
I have provided a base public directory [here](https://github.com/mongolar/public_directory_example).

Clone the repository wherever you wish to edit your html, js and css for the website.

I am just going to provide instructions on bower installation

Inside the directory run

```bash
bower install mongolar-js
```
This will install all the js libraries required in your index.html in a folder called assets.

The only file that will be accessible from the root folder of this site is the index.html.  Everything else must go in the assets folder.

If you want to change that functionality just go into the router and change the folder name to whatever you want.

###4. MongoDB
If you are new to MongoDB the easiest way to setup MongoDb is with [Mongolabs](https://mongolab.com/).  They offer a free tier service that will allow you to test the system, but don't expect the free tier service to be performant.  Latency between requests and the fact that the free tier does not have many resources can impact site performance.

Go to their site, create a database, and add a user to the database and they will provide the information to login from your Mongolar site.

You can seed the database from an export provided here LINK HERE

###5. Finalize
Add your MongoDB credentials and the site root directory to the configs file and boot mongolar.

You may need to run as root, depending on port and permissions of your OS.
```bash
sudo -E go run mongolar.go
```
You should be running, let me know if you have problems in the issue que.


##This is a very early BETA
This is in no way production ready.  There is still a lot to be done.

##Admin and OAuth controllers
I wrote two packages that are included with this code repository called admin and oauth, they are not well written and I rushed through them.
I wanted to create a UI where people could understand what this project does.  Those controller packages should not be considered production ready, 
and may not even be developed further.
If they do get developed further (read as severely overhauled), you can expect them to eventually be broken out to separate projects.

##Credits
There are several credits needed to be doled out.

First and foremost is my amazing wife.  She has really been patient with my efforts here.  She is really the best thing that has happened to me.

The Angular community.  They have to be the most supportive group of developers I have ever worked with.

The below package/library providers:

[formly-js](https://github.com/formly-js/angular-formly) - By far the best form generation tool I have ever used.

[davecgh/go-spew](https://github.com/davecgh/go-spew) - This package was a livesaver

[mgo](https://labix.org/mgo) -  Could not have done it without a MongoDb driver

[spf13/viper](https://github.com/spf13/viper)

[Sirupsen/logrus](https://github.com/Sirupsen/logrus)

##More information
Visit the issue que or read the Wiki.

##Want to help?
Create an issue

##Roadmap
  - Code Cleanup
  - Tests
  - Mock wrapper for testing
  - Kahn: a cli to talk to your mongolar server while it is running.
  - Clustering

##This sounds complicated

It is and it isn't.
