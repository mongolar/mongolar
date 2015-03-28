var mongular = angular.module('mongular', ['formly', 'angular-growl', 'ui.bootstrap', 'yaru22.md', 'ngSanitize']);

mongular.config(["$httpProvider", function ($httpProvider) {
    $httpProvider.interceptors.push('messageCatcher');
    $httpProvider.interceptors.push('reloadCatcher');
    $httpProvider.interceptors.push('redirectCatcher');
    $httpProvider.interceptors.push('queryCatcher');
    $httpProvider.interceptors.push('dynamicCatcher');
}]);
mongular.run(function($http) {
    $http.defaults.headers.common.CurrentPath = window.location.pathname;
});

mongular.controller('ContentController', function ContentController($scope, mongularService) {
    mongularService.mongularHttp($scope).then(function(response) {
        if(typeof response.data =='object') {
            angular.extend($scope, response.data);
        }
    });
    $scope.onSubmit = function() {
        var form_data = {};
        console.log($scope);
        $scope.mongular_content.formFields.forEach(function(field){
            form_data[field.key] = $scope.mongular_content.formData[field.key];
        });
        form_data.form_id = $scope.mongular_content.formOptions.uniqueFormId;
        mongularService.mongularPost(form_data).then(function(response) {
            if(typeof response.data =='object') {
                angular.extend($scope, response.data);
            }
        });
    }
    $scope.dynLoad = function(){
        mongularService.mongularHttp($scope).then(function(response) {
            if(typeof response.data =='object') {
                angular.extend($scope, response.data);
            }
        });
    }
    $scope.scopeSend =  function($path, $value) {
        mongularService.mongularScopeSend($path, $value);
    }
    $scope.serverSend = mongularService.mongularServerSend;
});
mongular.factory('mongularService', function($http, growl, mongularConfig) {
    var mongularService = {
        mongularHttp: function($arguments) {
            var argument = $arguments.mongulartype;
            if($arguments.mongular != undefined){
                argument = argument + '/' + $arguments.mongular;
            }
            var promise = $http.get(mongularConfig.mongular_url + argument).success(function(response){
                return response;
            });
            return promise;
        },
        mongularPost: function($arguments) {
            var promise = $http.post(mongularConfig.mongular_url + "form", $arguments).success(function(response){
                return response;
            });
            return promise;
        },
        mongularServerSend: function($path) {
            var promise = $http.get(mongularConfig.mongular_url  + $path).success(function(response){
                return response;
            });
            return promise;
        },
        mongularScopeSend: function($path, $values) {
            var promise = $http.post(mongularConfig.mongular_url  + $path, $values).success(function(response){
                return response;
            });
            return promise;
        }
    };
    return mongularService;
});

mongular.factory('queryCatcher', function($injector) {
    var queryCatcher = {
        response: function(response) {
            if(typeof response.data.mongular_content =='object') {
                if('query_parameters' in response.data.mongular_content){
                    var http = $injector.get('$http');
                    http.defaults.headers.common.QueryParameters = response.data.mongular_content.query_parameters;
                }
            }
            return response;
        }
    }
    return queryCatcher;
});

mongular.factory('messageCatcher', function(growl) {
    var messageCatcher = {
        response: function(response) {
            if(typeof response.data =='object') {
                if('mongular_messages' in response.data) {
                    response.data.mongular_messages.forEach(function (message) {
                        var growlcall = 'add' + message.severity + 'Message';
                        growl[growlcall](message.text);
                    });
                }
            }
            return response;
        }
    }
    return messageCatcher;
});

mongular.factory('reloadCatcher', function() {
    var reloadCatcher = {
        response: function(response) {
            if(typeof response.data =='object') {
                if('mongular_reload' in response.data){
                    location.reload();
                }
            }
            return response;
        }
    }
    return reloadCatcher;
});

mongular.factory('redirectCatcher', function() {
    var reloadCatcher = {
        response: function(response) {
            if(typeof response.data =='object') {
                if('mongular_redirect' in response.data){
                    window.location.replace(response.data.mongular_redirect);
                }
            }
            return response;
        }
    }
    return reloadCatcher;
});

mongular.factory('dynamicCatcher', function($rootScope) {
    var dynamicCatcher = {
        response: function(response) {
            if(typeof response.data =='object') {
                if('mongular_dynamics' in response.data){
                    response.data.mongular_dynamics.forEach(function (dynamic) {
                        $rootScope.$broadcast(dynamic.target, {
                            dyn_control: dynamic.controller,
                            dyn_template: dynamic.template,
                            dyn_id: dynamic.id,
                        });
                    });
                }
            }
            return response;
        }
    }
    return dynamicCatcher;
});

mongular.directive('mongularWrapper', function mongularWrapper($http, $compile, $rootScope, mongularConfig){
  return {
    controller: 'ContentController',
    transclude: true,
    scope: {
      'mongular' : '@',
      'mongulartemplate' : '@',
      'mongulartype' : '@',
      'mongulardyn' : '@'
    },
    template: '<div ng-include = "getTemplate()" ng-hide="hide"></div>',
      link : function(scope)
      {
          scope.getTemplate = function(){
              return mongularConfig.templates_url + scope.mongulartemplate;
          }
          if (scope.mongulardyn != undefined){
              scope.$on(
                  scope.mongulardyn, function (event, data) {
                      scope.mongulartype = data.dyn_control;
                      scope.mongular = data.dyn_id;
                      scope.mongulartemplate = data.dyn_template;
                      scope.dynLoad();
                  }
              );
          }
          scope.dynamicLoad = function(mongular_id, controller, template, target) {
              $rootScope.$broadcast(target, {
                  dyn_control : controller,
                  dyn_template : template,
                  dyn_id : mongular_id
              });
          }
      }
  }
});

