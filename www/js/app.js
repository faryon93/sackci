// ------------------------------------------------------------------------------------------------
var app = angular.module('app', ["ngRoute", "ngResource"]);
app.config(function($routeProvider, $locationProvider) {
    $routeProvider
        .when("/", {
            templateUrl : "project.html"
        })
        .when("/admin", {
            templateUrl : "admin.html"
        })
        .when("/project/:id", {
            templateUrl : "project.html"
        })
        .when("/project/:id/:tab", {
            templateUrl : "project.html"
        })
        .when("/project/:id/:tab/:build", {
            templateUrl : "project.html"
        });

    $locationProvider.html5Mode(true);
});


// ------------------------------------------------------------------------------------------------
app.controller("navigation", function($scope, $location) {
    $scope.isActive = function (viewLocation) { 
        return $location.path().startsWith(viewLocation);
    };
});

app.controller("projectlist", function($scope, $location, projects, feed) {
    // load project list from backend
    $scope.loading = true;
    $scope.projects = projects.query(function() {
        $scope.loading = false;

        // we are interested in all project status changes
        feed.register("project_changed", $scope, function(evt) {
            angular.forEach($scope.projects, function(project) {
                if (project.id === evt.project)
                    project.status = evt.status;
            });
        });       
    });
});

app.controller("project", function($scope, $location, $routeParams, projects, trigger) {
    // default tab is the history tab
    $scope.tab = $routeParams.tab;
    if ($scope.tab === undefined)
        $scope.tab = "history";

    // no project specified -> display empty page
    // TODO: display empty page
    if ($routeParams.id === undefined)
        return;

    $scope.build = $routeParams.build;

    // load the project details from the backend
    $scope.loaded = false;
    $scope.project = projects.get({id: $routeParams.id}, function() {
        $scope.loaded = true;
    });

    // event handler
    $scope.trigger = function() {
        var response = trigger.get({id: $routeParams.id}, function() {
            $location.path("/project/" + $routeParams.id + "/build/" + response.build_id);
        });
    }
});

app.controller("projectBuild", function($scope, $routeParams, builds) {
    $scope.select = function(stage) {
        $scope.stage = stage;
    };

    // redirect to latest ->
    var buildId = $routeParams.build;
    if (buildId === undefined)
        buildId = "latest";

    $scope.build = builds.get({id: buildId},
        function(response){
            $scope.stage = response.stages[0];
        },
        function(error) {
            $scope.error = error.data;
        }
    );
});

app.controller("projectHistory", function($scope, $routeParams, history) {
    $scope.builds = history.query({project: $routeParams.id},
        function(response) {},
        function(error) {
            $scope.error = error.data;
        });
});

app.controller("projectEnv", function($scope, $routeParams, env) {
    $scope.refresh = function() {
        // load environment information from the backend
        $scope.envs = env.query({id: $routeParams.id});
    };

    $scope.showConfirm = function(show) {
        $scope.confirm = show;
    };

    // the initial refresh
    $scope.refresh();
});


// ------------------------------------------------------------------------------------------------
app.factory('projects', function($resource) {
    return $resource("/api/v1/project/:id");
});

app.factory('builds', function($resource){
   return $resource('/api/v1/build/:id');
});

app.factory('history', function($resource){
   return $resource('/api/v1/project/:project/history/:id', {project:'@project', id: '@id'});
});

app.factory('env', function($resource){
   return $resource('/api/v1/project/:id/env');
});

app.factory('trigger', function($resource){
    return $resource('/api/v1/project/:id/trigger');
});

// ------------------------------------------------------------------------------------------------
app.factory('feed', function($rootScope) {
    var sse = new EventSource('/api/v1/feed');
    console.log("feed created");

    return {
        register: function(eventName, scope, callback) {
            // register the event listener
            var fn = function(evt) {
                $rootScope.$apply(function () {
                    callback.apply(sse, [JSON.parse(evt.data)]);
                });
            };
            sse.addEventListener(eventName, fn);

            // remove the event listener, wenn the calling scope
            // is destroyed in order to unregister from the event source
            scope.$on('$destroy', function() {
                sse.removeEventListener(eventName, fn);
            });
        }
    };
});


// ------------------------------------------------------------------------------------------------
app.filter("timeago", function () {
    //time: the time
    //local: compared to what time? default: now
    //raw: wheter you want in a format of "5 minutes ago", or "5 minutes"
    return function (time, local, raw) {
        if (!time || time === "0001-01-01T00:00:00Z")
            return "never";

        if (!local) {
            (local = Date.now())
        }

        if (angular.isDate(time)) {
            time = time.getTime();
        } else if (typeof time === "string") {
            time = new Date(time).getTime();
        }

        if (angular.isDate(local)) {
            local = local.getTime();
        }else if (typeof local === "string") {
            local = new Date(local).getTime();
        }

        if (typeof time !== 'number' || typeof local !== 'number') {
            return;
        }

        var
            offset = Math.abs((local - time) / 1000),
            span = [],
            MINUTE = 60,
            HOUR = 3600,
            DAY = 86400,
            WEEK = 604800,
            MONTH = 2629744,
            YEAR = 31556926,
            DECADE = 315569260;

        if (offset <= MINUTE)              span = [ '', raw ? 'now' : 'less than a minute' ];
        else if (offset < (MINUTE * 60))   span = [ Math.round(Math.abs(offset / MINUTE)), 'min' ];
        else if (offset < (HOUR * 24))     span = [ Math.round(Math.abs(offset / HOUR)), 'hr' ];
        else if (offset < (DAY * 7))       span = [ Math.round(Math.abs(offset / DAY)), 'day' ];
        else if (offset < (WEEK * 52))     span = [ Math.round(Math.abs(offset / WEEK)), 'week' ];
        else if (offset < (YEAR * 10))     span = [ Math.round(Math.abs(offset / YEAR)), 'year' ];
        else if (offset < (DECADE * 100))  span = [ Math.round(Math.abs(offset / DECADE)), 'decade' ];
        else                               span = [ '', 'a long time' ];

        span[1] += (span[0] === 0 || span[0] > 1) ? 's' : '';
        span = span.join(' ');

        if (raw === true) {
            return span;
        }
        return (time <= local) ? span + ' ago' : 'in ' + span;
    }
});

app.directive('capitalize', function() {
    return {
      require: 'ngModel',
      link: function(scope, element, attrs, modelCtrl) {
        var capitalize = function(inputValue) {
          if (inputValue === undefined)
              inputValue = '';

          var capitalized = inputValue.toUpperCase();
          if (capitalized !== inputValue) {
              modelCtrl.$setViewValue(capitalized);
              modelCtrl.$render();
          }

          return capitalized;
        };

        modelCtrl.$parsers.push(capitalize);
        capitalize(scope[attrs.ngModel]);
      }
    };
});

app.filter("ansicolor", ['$sce', function($sce) {
    var foregroundColors = {
      '01': 'bold',
      '30': 'black',
      '31': 'red',
      '32': 'green',
      '33': 'yellow',
      '34': 'blue',
      '35': 'magenta',
      '36': 'cyan',
      '37': 'white'
    };

    return function(val)
    {
        if (val === undefined)
            val = '';

        Object.keys(foregroundColors).forEach(function (ansi) {
            var span = '<span class="ansi ' + foregroundColors[ansi] + '">';
            var boldSpan = '<span class="ansi bold ' + foregroundColors[ansi] + '">';
            val = val.replace(new RegExp('\033\\[' + ansi + 'm', 'g'), span)
                     .replace(new RegExp('\033\\[0;' + ansi + 'm', 'g'), span)
                     .replace(new RegExp('\033\\[01;' + ansi + 'm', 'g'), boldSpan);
        });

        // bold and italic
        val = val.replace(/\033\[1m/g, '<b>').replace(/\033\[22m/g, '</b>');
        val = val.replace(/\033\[3m/g, '<i>').replace(/\033\[23m/g, '</i>');

        // closing tag
        val = val.replace(/\033\[m/g, '</span>');
        val = val.replace(/\033\[0m/g, '</span>');
        val = val.replace(/\033\[39m/g, '</span>');

        return $sce.trustAsHtml(val);
    };
}]);

app.filter('duration', function(){
    return function(input){
        // the time is supplied in microseconds
        input = input / 1000.0 / 1000.0 / 1000.0;

        var minutes = parseInt(input/60, 10);
        var seconds = Math.ceil(input % 60);

        var result = "";
        if (minutes > 0)
            result += minutes + "min ";

        if (seconds > 0)
            result += seconds + "sec";

        return result;
    }
});

app.filter('reverse', function() {
    return function(items) {
        return items.slice().reverse();
    };
});

app.filter('unknown', ['$sce', function($sce) {
    return function(str) {
        var style = "";
        if (str === "unknown")
            style += "font-style: italic;";

        return $sce.trustAsHtml('<span style="' + style + '">' + str + '</span>');
    };
}]);

app.filter('short', function() {
    return function(str) {
        if (str === undefined)
            return "";

        return str.substr(0, 12);
    }
});