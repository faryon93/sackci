var app = angular.module('app', ["ngRoute", "ngResource"]);

app.config(function($routeProvider) {
    $routeProvider
    .when("/", {
        templateUrl : "project.html"
    })
    .when("/project/:id", {
        templateUrl : "project.html"
    })
    .when("/project/:id/:tab", {
        templateUrl : "project.html"
    })
    .when("/project/:id/build/:build", {
        templateUrl : "project.html"
    });
});

app.controller("navigation", function($scope, $location) {
	$scope.isActive = function (viewLocation) { 
        return $location.path().startsWith(viewLocation);
    };
});

app.controller("projectlist", function($scope, $location, projects) {
	$scope.projects = projects.get()
});

app.controller("project", function($scope, $location, $routeParams, projects) {
	$scope.project = projects.find($routeParams.id == undefined ? 0 : $routeParams.id);
	$scope.tab = $routeParams.tab;
	if ($scope.tab == undefined)
		$scope.tab = "build";

	$scope.build = $routeParams.build;
});

app.controller("projectenv", function($scope) {
    console.log($scope.project);
});

app.factory('projects', function() {
	var projects = [{
        id: 0,
        name: "org.kde.breeze",
        build: 11,
        execution_time: "9 days ago",
        duration: "1 minute",
        passed: true,
        env: [],
    },
    {
        id: 1,
        name: "tcdeploy",
        build: 5,
        execution_time: "5 hours ago",
        duration: "12 seconds",
        passed: true,
        env: [],
    },
    {
        id: 2,
        name: "Xorbit 0.8 Legacy",
        build: 432,
        execution_time: "29 minutes ago",
        duration: "7 minutes",
        passed: false,
        env: [
        	{key: "XORBIT_BUILD_TIME", value: "$now()", encrypted: false},
        	{key: "XORBIT_DB_HOST", value: "172.16.1.59", encrypted: false},
        	{key: "XORBIT_DB_USER", value: "root", encrypted: false},
        	{key: "XORBIT_DB_PASSWORD", value: "", encrypted: true}
        ],
    }];

	var service = {};
    service.get = function() {
        return projects;
    };
    service.find = function(id) {
    	p = projects[id];
    	return p;
    };
    
    // other stubbed methods
    return service;
});