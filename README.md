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

Every piece of page content (read API request) is individually addressable so if situated behind a tool like varnish, you can pick and choose caching behavior.

MongoDB seems to scale rather well, at least for these purposes.

## Demo
I plan to release a demo for this server shortly.  In the demo you will be able to manipulate content and view the results

##Setup

To setup the go server you will do several things (assuming go is already setup on your system):

###1. Download mongolar.
```bash
go get github.com/mongolar/mongolar
```

###2. Configuration
Mongolar currently only has one configuration file for a server and as many site files as you wish

All configurations are in YAML format

The default location for a server config would be here "/etc/mongolar", if this does not work for you you can set a "MONGOLAR_SERVER_CONFIG" environment variable to the folder you want to use.

Example:
```yaml
"Port": "80"
#Directory where site configs will be stored
"SitesDirectory": "/etc/mongolar/enabled/"
#Directory for logs
"LogDirectory": "/var/log/mongolar/"
```

Based on the server config Mongolar will attempt to load all site configuration files into memory, so given the above configuration you would create a yaml:

"/etc/mongolar/enabled/my_site.yaml" <--yaml suffix required

Ideally you would create your configs in /etc/mongolar/available/my_site.yaml and then symlink to the file, making configs easier to manage.

Example:
```yaml
# The Mongodb connection
# REQUIRED
MongoDb:
        user: my_db_user
        password: my_password
        host: "db_domain:12345"
        db: my_db
# Where you will store site configurations
# REQUIRED
Directory: "/my/files/directory"
# All the domains that will apply to this site
# REQUIRED
Aliases:
        - mongolar.org
        - test-1.com
        - test-1.1.com
# Public values are values that apply site wide, 
# and can safely be served stright from the config
# An example may be your google analytics code
PublicValues:
        "test" : "Test Value"
# Page not found, you can set this to whatever you want.
# Requires the page be built for the url
# REQUIRED
FourOFour: "/page_not_found"
# Per site angular module loading, this will be added to the mongolar config
# REQUIRED
AngularModules:
        - "formly" #Required for forms
        - "formlyBootstrap" #Other formly templates available
        - "angular-growl" #
        - "ui.bootstrap"
        - "ngSanitize"
        - "ui.sortable"
# Where templates will originate from, this can be a cdn
# REQUIRED
TemplateEndpoint: "assets/templates"
# Foreign domains allows AngularJS to load assets from domains outside your own
# This is required if you set your templates from a cdn
ForeignDomains: "my.foo.com"
# When to expire Session after so many hours
# REQUIRED
SessionExpiration: 10
# The location where your API can be reached in your domain, can be any string
# REQUIRED
APIEndPoint: "my_end_point"
# Each Website must specify  which controllers it can utilize.
# Even if the controller is compiled in the binary,
# if it is not listed here it will be forbidden.
# This allows you to restrict access to controllers per site.
# REQUIRED
Controllers:
        - "path"
        - "content"
        - "wrapper"
        - "slug"
        - "admin"
        - "domain_public_value"
        - "login"
# For the current incarnation of Mongolar this works,
# but will most likely be changed
# Stores OAuthlogins, currently only supports github
OAuthLogins:
        "github":
                "client_id": "client_id_here"
                "client_secret": "client_secret_here"
                "login_text": "Login with github"
# This is a standard set of urls where you can expect to send
# users for login functions etc.
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

You can seed the database from an export provided [here](https://github.com/mongolar/seed_database).

To import you will need to use mongorestore, check with your OS and how to install.

```bash
mongorestore -h the_endpoint_mongolabs provided -d database_name -u user_name -p password the_directory_you_just_cloned
```

###5. Finalize
Add your MongoDB credentials and the site root directory to the configs file and boot mongolar.

You may need to run as root, depending on port and permissions of your OS.
```bash
sudo -E mongolar
```
You should be running, let me know if you have problems in the issue que.

###5.1 Access Admin UI

If you are using the seed DB and want to access the admin ui, you will have to do the following.
Login using the your github credentials under "/login"
Now under mongolabs, go to your db, under the "users" collection.  Your login should be the only document.
Edit that record and add the following to the root of your document
```json
"roles": [
	"admin"
]
```
Make sure you have the correct commas after the value above the roles value.

Example from mine:
```json
{
    "_id": {
        "$oid": "the object id"
    },
    "email": "my email address@gmail.com",
    "id": 1111111,
    "name": "jasonrichardsmith",
    "type": "github",
    "roles": [
        "admin"
    ]
}
```
Values were changed to protect the innocent.

##This is a very early BETA
This is in no way production ready.  There is still a lot to be done.


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

This list will grow for sure.

##More information
Visit the issue que or read the Wiki.

##Want to help?
Fork it!

##Feature Roadmap
  - Code Cleanup
  - Tests
  - Mock wrapper for testing
  - Kahn: a cli to talk to your mongolar server while it is running.
  - Clustering

##This sounds complicated

It is and it isn't.

##Disclaimer

I learned Go writing this project, and I got this far in two months of my free time (which I have little) so there are some issues...
Here is a list things I know need to be fixed.
- Element loading from mgo,  It needs to accept structures to avoid interfaces and reflection.
- Post data, it is being parsed by the wrapper and this responsibility should be passed to the controllers, so controllers can marshall their own expectations.
- Logic to test for Post should be in the controller and should be based off the method
- Sessions are upserting on every api request (which is atomic), this needs to be improved (would love input on this one)
- Packages need to be broken into smaller files.
- The API URL parsing is incorrect but works
- URL parameters need to be validated to avoid panics in the controllers
- ObjectIds need to be validated prior to being added to bson.M

This list goes on.

##Slug values and Wildcard paths
The system does support wildcard paths which means, if an explicit path does not match it will attempt to retrieve the wildcard path.
Angular then appends the rest of the url (the part after the wildcard match) to the header of each subsequent request.
This means you can have one a "/blog" path that loads the same way each time but loads data based on a slug value.
You can see the [slug controller](https://github.com/mongolar/mongolar/blob/master/controller/controller.go#L242) on how this is achieved.

I have not built anything in the Admin UI to administer this.


##Admin and OAuth controllers
I wrote two packages that are included with this code repository called admin and oauth, they are not well written and I rushed through them.
I wanted to create a UI where people could understand what this project does.  Those controller packages should not be considered production ready, 
and may not even be developed further.

If they do get developed further (read as severely overhauled), you can expect them to eventually be broken out to separate projects.
