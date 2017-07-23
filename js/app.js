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
});


// ------------------------------------------------------------------------------------------------
app.controller("navigation", function($scope, $location) {
    $scope.isActive = function (viewLocation) { 
        return $location.path().startsWith(viewLocation);
    };
});

app.controller("projectlist", function($scope, $location, projects) {
    // load project list from backend
    $scope.loading = true;
    $scope.projects = projects.query(function() {
        $scope.loading = false;
    });
});

app.controller("project", function($scope, $location, $routeParams, projects) {
    // default tab is the history tab
    $scope.tab = $routeParams.tab;
    if ($scope.tab == undefined)
        $scope.tab = "history";

    // no project specified -> display empty page
    // TODO: display empty page
    if ($routeParams.id == undefined)
        return;

    $scope.build = $routeParams.build;

    // load the project details from the backend
    $scope.loaded = false;
    $scope.project = projects.get({id: $routeParams.id}, function() {
        $scope.loaded = true;
    });
});

app.controller("projectBuild", function($scope, $routeParams, builds) {
    // redirect to latest ->
    var buildId = $routeParams.build;
    if (buildId == undefined)
        buildId = "latest";

    $scope.build = builds.get({project: $routeParams.id, id: buildId});
});

app.controller("projectHistory", function($scope, $routeParams, builds) {
    $scope.builds = builds.query({project: $routeParams.id});
});

app.controller("projectEnv", function($scope, $routeParams, env) {
    $scope.question = function(item, show) {
        item.question = show;
    };

    $scope.refresh = function() {
        // load environment information from the backend
        $scope.envs = env.query({project: $routeParams.id});
    };

    // the initial refresh
    $scope.refresh();
});

// ------------------------------------------------------------------------------------------------
app.factory('projects', function($resource) {
    return $resource("/api/v1/project/:id");
});

app.factory('builds', function($resource){
   return $resource('/api/v1/project/:project/build/:id', {project:'@project', id: '@id'})
});

app.factory('env', function($resource){
   return $resource('/api/v1/project/:project/env/:id', {project:'@project', id: '@id'})
});


// ------------------------------------------------------------------------------------------------
app.filter("timeago", function () {
    //time: the time
    //local: compared to what time? default: now
    //raw: wheter you want in a format of "5 minutes ago", or "5 minutes"
    return function (time, local, raw) {
        if (!time) return "never";

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
          if (inputValue == undefined) inputValue = '';
          var capitalized = inputValue.toUpperCase();
          if (capitalized !== inputValue) {
            modelCtrl.$setViewValue(capitalized);
            modelCtrl.$render();
          }
          return capitalized;
        }

        modelCtrl.$parsers.push(capitalize);
        capitalize(scope[attrs.ngModel]);
      }
    };
});