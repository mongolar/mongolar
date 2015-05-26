#Mongolar

Mongolar is a AngularJS driven scalable Content Management System written in AngularJS and Go with MongoDB as a backend.

##Basics
###AngularJS
Essentially AngularJS bootstraps itself upon page load.  It then begins to process any "mongolar" directives, loading content to the page.

Mongolar directives in AngularJS perform api calls based on settings in the directive that give the address for the call, an id, a template and a dynamic id (an element id that makes dynamic page content loading easier).

The simplest implementation for a directive would be to load page elements from a path and then the returned values would then create more directives to subseqquently make more API requests

###Go

The Mongolar binary will work as a webserver/API server.

Mongolar serves along several ddifferent paths
 - /assets : Holds all html assets
 - /mongolarconfig.js : This is a site specific js file that loads all the mongolar-js specific configurations for a particular site
 - /apiendpoint : This a per site configurable endpoint for api requests.  This can be any url friendly string
 - default: If all the above routing does not work then the index.html will be served allowing AngularJS to take over routing.

###MongoDB
This system was written with heavy reads in mind.  Most web content will be heavy on the reads, light on the writes.
I also wanted a NoSQL Db that was easy to setup, learn and scale.

##This sounds complicated

It is and it isn't.

##Scalable

This package seeks to achieve scalability in several ways.

Each Mongolar server instance is stateless.  You can spin up as many servers as you wish behind a load balancer, and there is not extra configuration.

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
EXAMPLE HERE
```
###3. Frontend
I have provided a base public directory here. LINK HERE

Clone the repository wherever you wish to edit your html, js and css for the wwebsite.

I am just going to provide instructions on bower installation

Inside the directory run

```bash
bower install mongolar-js
```
This will install all the js libraries required in your index.html

The only file that will be accessible from the root folder of this site is the index.html.  Everything else must go in the assets folder.

If you want to change that functionality just go into the router and change the folder name to whatever you want.

###4. MongoDB
If you are new to MongoDB the easiest way to setup MongoDb is with Mongolabs.

Go to their site, create a database, and add a user to the database and they will provide the information to login from your Mongolar site.

You can seed the database from an export provided here LINK HERE

###5. Finalize
Add your MongoDB credentials and the site root directory to the configs file and boot mongolar.

You may need to run as root, depending on port and permissions of your OS.
```bash
sudo -E go run mongolar.go
```

