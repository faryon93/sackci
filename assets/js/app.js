// ------------------------------------------------------------------------------------------------
var app = angular.module('app', ["ngRoute", "ngResource", "ui.bootstrap", "infinite-scroll"]);
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
        feed.register("EvtPipelineBegin", $scope, function(evt) {
            angular.forEach($scope.projects, function(project) {
                if (project.id === evt.project_id)
                {
                    project.status = evt.status;
                    project.execution_time = evt.time;
                    project.build_num = evt.build_num;
                    project.duration = 0;
                }
            });
        });
        feed.register("EvtPipelineFinished", $scope, function(evt) {
            angular.forEach($scope.projects, function(project) {
                if (project.id === evt.project_id)
                {
                    project.status = evt.status;
                    project.duration = evt.duration;

                    if (evt.status === "waiting")
                    {
                        project.build_num = 0;
                        project.execution_time = 0;
                    }
                }
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

    // event handler: trigger build
    $scope.trigger = function() {
        $scope.triggerStatus = "waiting";
        $scope.triggerError = "";

        var response = trigger.get({id: $routeParams.id}, function() {
            $location.path("/project/" + $routeParams.id + "/build/" + response.build_id);
        }, function (error) {
            $scope.triggerStatus = "error";
            $scope.triggerError = error.data;
        });
    };

    $scope.triggerBlur = function() {
        $scope.triggerStatus = "";
        $scope.triggerError = "";
    };
});

app.controller("projectBuild", function($scope, $routeParams, builds, feed) {
    // default variables
    $scope.timestamp = 0;

    // handler: select a stage
    $scope.selectStage = function(stage) {
        if (stage.status !== "ignored")
            $scope.stage = stage;
    };

    // successfull loaded build details
    var success = function(response, headers) {
        // get the serverside timestamp in order to discard old feed events
        $scope.timestamp = headers()["x-timestamp"];

        // per default the stage wich is most recent
        // and not ignored is displayed when viewing the build
        $scope.build.stages.forEach(function(stage) {
            if (stage.status !== "ignored")
                $scope.stage = stage;
        });
    };

    // error while loading build details
    var error = function(error) {
        $scope.error = error.data;
        $scope.errorCode = error.status;
    };

    // register for an event from the newsfeed
    var registerFeed = function(event, callback) {
        feed.register(event, $scope, function(evt) {
            // make sure the event is for the selected build
            // if not, discard the event
            if (($scope.timestamp > evt.timestamp) ||
                $scope.project.id !== evt.project_id ||
                $scope.build.num !== evt.build_num)
            {
                return;
            }

            // call the actual function
            if (typeof callback === "function")
                callback(evt);
        });
    };

    // if no proper build id is given ask for the latest
    var buildId = $routeParams.build;
    if (buildId === undefined)
        buildId = "latest";
    $scope.build = builds.get({project: $routeParams.id, id: buildId}, success, error);

    // register for updates from the news feed
    registerFeed("EvtCommitFound", function(evt) {
        $scope.build.commit = evt.commit;
    });
    registerFeed("EvtPipelineFinished", function(evt) {
        $scope.build.status = evt.status;
        $scope.build.duration = evt.duration;
    });
    registerFeed("EvtStageBegin", function(evt) {
        $scope.build.stages[evt.stage].status = evt.status;
        $scope.stage = $scope.build.stages[evt.stage];
    });
    registerFeed("EvtPipelineFound", function(evt) {
        evt.stages.forEach(function(stageName) {
            // TODO: the server should send the full stage structure
            $scope.build.stages.push({name: stageName, status: "ignored", duration: 0, log: []});
        });
    });
    registerFeed("EvtStageLog", function(evt) {
        var stageId = evt.stage;
        if ($scope.build.stages[stageId] !== undefined)
            $scope.build.stages[stageId].log.push(evt.message);
    });
    registerFeed("EvtStageFinish", function(evt) {
        var stageId = evt.stage;
        if ($scope.build.stages[stageId] !== undefined)
        {
            $scope.build.stages[stageId].status = evt.status;
            $scope.build.stages[stageId].duration = evt.duration;
        }
    });
});

app.controller("projectHistory", function($scope, $routeParams, history, feed) {
    var projectId = parseInt($routeParams.id);
    var limit = 10;
    var skip = 0;

    // eventhandler: infinite page scroll
    $scope.scrollBusy = true;
    $scope.scrollEnd = false;
    $scope.nextPage = function() {
        $scope.scrollBusy = true;

        // get the next page from the api
        skip += limit;
        var builds = history.query({id: $routeParams.id, limit: limit, skip: skip},
            function() {
                // if the server sends an empty array no new
                // elements are available
                if (builds.length === 0)
                    $scope.scrollEnd = true;

                // append the new page of builds to the other ones
                $scope.builds = $scope.builds.concat(builds);
                $scope.scrollBusy = false;
            });
    };

    // query for the build history
    $scope.builds = history.query({id: $routeParams.id, limit: limit, skip: skip},
        function() {
            $scope.scrollBusy = false;
            if ($scope.builds.length < limit)
                $scope.scrollEnd = true;
        },
        function(error) {
            $scope.error = error.data;
        });

    // register for some events
    feed.register("EvtPipelineBegin", $scope, function(evt) {
        if (evt.project_id !== projectId)
            return;

        var found = false;
        angular.forEach($scope.builds, function(build) {
            if (build.num === evt.build_num)
            {
                found = true;
                build.status = evt.status;
                build.time = evt.time;
                build.node = evt.agent;
                build.duration = 0;
            }
        });

        // the build is not in the liste -> create a new one and append
        if (!found)
        {
            $scope.builds.unshift({
                num: evt.build_num,
                status: evt.status,
                node: evt.agent,
                time: evt.time,
                duration: 0
            })
        }
    });
    feed.register("EvtPipelineFinished", $scope, function(evt) {
        if (evt.project_id !== projectId)
            return;

        // a build num 0 means the build history
        // has been purged
        if (evt.build_num === 0)
        {
            $scope.builds = [];
            return;
        }

        angular.forEach($scope.builds, function(build) {
            if (build.num === evt.build_num)
            {
                build.status = evt.status;
                build.duration = evt.duration;
            }
        });
    });
    feed.register("EvtCommitFound", $scope, function(evt) {
        if (evt.project_id !== projectId)
            return;

        angular.forEach($scope.builds, function(build) {
            if (build.num === evt.build_num)
                build.commit = evt.commit;
        });
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

app.controller("settings", function($scope, $routeParams, history) {
    // eventhandler: purge build history
    $scope.purge = function(scope) {
        scope.status = "waiting";
        history.remove({id: $routeParams.id}, function() {
                scope.status = "success";
                scope.error = ""
            },
            function (error) {
                scope.status = "error";
                scope.error = error.data;
            });
    };

    // eventhandler: delete project
    $scope.delete = function(scope) {

    };
});


// ------------------------------------------------------------------------------------------------
app.factory('projects', function($resource) {
    return $resource("/api/v1/project/:id");
});

app.factory('builds', function($resource){
   return $resource('/api/v1/project/:project/build/:id', {project:'@project', id: '@id'});
});

app.factory('history', function($resource){
   return $resource('/api/v1/project/:id/history?limit=:limit&skip=:skip',
       {limit:'@limit', skip: '@skip', id: '@id'});
});

app.factory('env', function($resource){
   return $resource('/api/v1/project/:id/env');
});

app.factory('trigger', function($resource){
    return $resource('/api/v1/project/:id/trigger');
});


// ------------------------------------------------------------------------------------------------
app.factory('feed', function($rootScope) {
    $rootScope.liveUpdates = true;
    var sse = new EventSource('/api/v1/feed');
    sse.onerror = function() {
        // TODO: retry connection
        $rootScope.$apply($rootScope.liveUpdates = false);
    };
    sse.onmessage = function() {
        $rootScope.$apply($rootScope.liveUpdates = true);
    };

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

        if (offset <= 10)
        {
            span = ['', "just now"];
            raw = true;
        }
        else if (offset <= MINUTE)         span = [ '', raw ? 'now' : 'a minute' ];
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

        return str.substr(0, 7);
    }
});

app.filter('badgeUrl', function() {
    return function(project) {
        if (project === undefined)
            return "";

        return window.location.origin + "/api/v1/project/" + project.id + "/badge";
    }
});


// ------------------------------------------------------------------------------------------------
app.directive("confirm", function ($compile) {
    return {
        scope: {
            confirm: '&'
        },
        link: function (scope, element, attrs) {
            attrs.$set('ng-click', 'showConfirmation()');
            attrs.$set('ng-show', '!show');
            element.addClass('confirm');
            element.removeAttr("confirm");

            var contentTr = angular.element(
                '<strong ng-show="show">Sure?</strong>&nbsp;' +
                '<button type="button" class="btn btn-default" ng-show="show" ng-click="yes()">' +
                '<i class="fa fa-check"></i> Yes</button>&nbsp;' +
                '<button type="button" class="btn btn-default" ng-show="show" ng-click="no()">' +
                '<i class="fa fa-remove"></i> No</button>');
            element.after($compile(contentTr)(scope));

            $compile(element)(scope);

            scope.showConfirmation = function() {
                scope.show = true;
                scope.test = '';
            };

            scope.yes = function() {
                scope.confirm()(scope);
                scope.show = false;
            };

            scope.no = function() {
                scope.show = false;
            }
        }
    }
});