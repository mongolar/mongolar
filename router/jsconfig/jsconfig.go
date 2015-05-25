package jsconfig

import (
	"net/http"
	"text/template"
)

// We are compiling the AngularJs config script into the program for fast loading.  All sites should have the same Angular JS config
// template with different values passed from the site config
const (
	ConfigScript = `
		var mongolar = angular.module('mongolar', 
			[
				{{ range .AngularModules }}
     					'{{.}}',
				{{ end }}
			]
		);
		
		mongolar.constant('mongolarConfig', {
    			mongolar_url: '/{{ .APIEndPoint }}/',
    			templates_url: '/{{ .TemplateEndpoint }}/'}
		);
		mongolar.config(function($sceDelegateProvider) {
   			$sceDelegateProvider.resourceUrlWhitelist([
     				// Allow same origin resource loads.
     				'self',
				{{ range .ForeignDomains }}
     				'http://{{.}}/**',
				{{ end }}
   			]);
 		});
		mongolar.config(['growlProvider', function(growlProvider) {
		  growlProvider.globalTimeToLive(5000);
		}]);
		`
)

// Current available values.  This may be reconfigured if it gets too large.
type JsConfigs struct {
	APIEndPoint      string
	TemplateEndpoint string
	ForeignDomains   []string
	AngularModules   []string
}

// Serve the config
func (c *JsConfigs) Serve(w http.ResponseWriter) {
	t := template.New("Mongolar Config JS")
	t.Parse(ConfigScript)
	w.Header().Set("Content-Type", "application/javascript")
	t.Execute(w, c)
}
