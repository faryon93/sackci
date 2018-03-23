// ------------------------------------------------------------------------------------------------
var app = angular.module('app', ["ngRoute", "ngResource", "ui.bootstrap", "infinite-scroll", "angular-loading-bar"]);
app.config(function($routeProvider, $locationProvider, $httpProvider) {
    $routeProvider
        .when("/", {
            templateUrl : "home.html"
        })
        .when("/login", {
            templateUrl: "login.html",
            fullpage: true
        })
        .when("/config", {
            templateUrl : "config.html",
            login: true
        })
        .when("/project/:id", {
            templateUrl : "project.html",
            login: true
        })
        .when("/project/:id/:tab", {
            templateUrl : "project.html",
            login: true
        })
        .when("/project/:id/:tab/:build", {
            templateUrl : "project.html",
            login: true
        });

    $locationProvider.html5Mode(true);
    $httpProvider.interceptors.push('authInterceptor')
});

// ------------------------------------------------------------------------------------------------
app.run(function($rootScope) {
    $rootScope.$on("$routeChangeStart", function(event, next, current) {
        if (next === undefined)
        {
            $rootScope.showProjectList = true;
            $rootScope.hideNavigation = false;
            return;
        }

        // TODO: check authentication state
        $rootScope.showProjectList = next.fullpage === undefined ? true : !next.fullpage;
        $rootScope.hideNavigation = next.fullpage === undefined ? false : next.fullpage;
    });
});

// ------------------------------------------------------------------------------------------------
app.controller("navigation", function($scope, $location, $http) {
    $scope.isActive = function (viewLocation) { 
        return $location.path().startsWith(viewLocation);
    };

    $scope.logout = function() {
        $http.get('/api/v1/logout');
    };
});

app.controller("login", function($scope, $http, $location) {
    $scope.form = {
        username: "",
        password: "",
        remember: false
    };

    $scope.login = function(scope) {
        $http.post("/api/v1/login", $scope.form)
            .then(function () {
                $location.path("/");
            },function (error){
                $scope.error = error.data;
                $scope.form.password = "";
                scope.reset();
            });

        return true;
    }
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

app.controller("project", function($scope, $location, $routeParams, $http, projects) {
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

        $http.post("/api/v1/project/" + $routeParams.id + "/trigger", "")
             .then(function(response) {
                 $location.path("/project/" + $routeParams.id + "/build/" + response.data.build_id);
             }, function(error) {
                 $scope.triggerStatus = "error";
                 $scope.triggerError = error.data;
             });

        // var response = trigger.save({id: $routeParams.id}, function() {
        //     $location.path("/project/" + $routeParams.id + "/build/" + response.build_id);
        // }, function (error) {
        //     $scope.triggerStatus = "error";
        //     $scope.triggerError = error.data;
        // });
    };

    $scope.triggerBlur = function() {
        $scope.triggerStatus = "";
        $scope.triggerError = "";
    };
});

app.controller("projectBuild", function($scope, $routeParams, builds, feed, log) {
    // constants
    const LOG_LIMIT = 50;

    // default variables
    $scope.timestamp = 0;
    $scope.errorCode = 404;
    $scope.error = "not found";
    $scope.loadingLog = false;
    $scope.expandCommit = false;

    var cleanTerminal = function(str) {

        // find the last line with actual text in it
        var lines = str.split("\r");
        for (var i = (lines.length - 1); i >= 0; i--)
        {
            const line = lines[i];
            if (line !== "" && line !== "\r" && line !== '\n' && line !== '\033[1B')
                return lines[i]
        }

        // fallback: nothing to split
        return str;
    };

    $scope.expand = function() {
        $scope.expandCommit = !$scope.expandCommit;
    };

    // handler: select a stage
    $scope.selectStage = function(index) {
        // do not display ignored stages
        var stage = $scope.build.stages[index];
        if (stage === undefined || stage.status === "ignored")
            return;

        if ($scope.stage === stage)
            return;

        // asign the dataobject to be rendered
        $scope.logLimit = LOG_LIMIT;
        $scope.stage = stage;

        // only display the loading spinner when no log
        // is displayed right now
        if ($scope.stage.log === undefined || $scope.stage.log.length === 0)
        {
            $scope.stage.log = [];
            $scope.loadingLog = true;
            // load the log of the stage
            log.get({project: $routeParams.id, id: $scope.build.num, stage: index},
                function(resp) {
                    $scope.stage.log = resp.log.split("\n").map(cleanTerminal);
                    $scope.loadingLog = false;
                }
                // TODO: error handling
            );
        }
    };

    // eventhandler: more logs lines
    $scope.moreLog = function() {
        $scope.logLimit += LOG_LIMIT;
    };

    // eventhandler: jump to end
    $scope.logJump = function () {
        $scope.logLimit = $scope.stage.log.length;
        $scope.$$postDigest(function() {
            window.scrollTo(0, document.body.scrollHeight);
        });
    };

    // successfull loaded build details
    var success = function(response, headers) {
        // get the serverside timestamp in order to discard old feed events
        $scope.timestamp = headers()["x-timestamp"];
        $scope.errorCode = 0;
        $scope.error = "";

        // per default the stage wich is most recent
        // and not ignored is displayed when viewing the build
        var selectedStage = 0;
        $scope.build.stages.forEach(function(stage) {
            if (stage.status !== "ignored")
                selectedStage++
        });

        if (selectedStage > 0)
            $scope.selectStage(selectedStage - 1);
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
        {
            $scope.logLimit++;
            $scope.build.stages[stageId].log.push(cleanTerminal(evt.message));
        }
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

app.controller("projectEnv", function($scope, $routeParams, projects) {
    $scope.addName = "";
    $scope.addValue = "";

    // eventhandler: add env variable button
    $scope.add = function() {
        $scope.showForm = true;
        $scope.addName = "";
        $scope.addValue = "";
    };

    // eventhandler: cancel adding of env variable
    $scope.cancelAdd = function() {
        $scope.showForm = false;
    };

    // eventhandler: save the new env variable
    $scope.save = function() {
        if ($scope.addName === "" || $scope.addValue === "")
            return;

        var obj = {};
        obj[$scope.addName] = $scope.addValue;
        projects.save({id: $routeParams.id}, {env: obj});

        $scope.showForm = false;
    };

    // delete an env variable
    $scope.delete = function (scope) {

    };
});

app.controller("settings", function($scope, $routeParams, $timeout, history, projects) {
    // variable initialisation
    $scope.viewKey = false;

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

    // eventhandler: rename project
    $scope.rename = function(scope) {
        if ($scope.changeName === undefined || $scope.changeName === "" ||
            $scope.project.name === $scope.changeName)
        {
            return false;
        }

        projects.save({id: $routeParams.id}, {name: $scope.changeName},
            function() {
                scope.success();
                $scope.project.name = $scope.changeName;
            },
            function(error) {
                scope.error(error.data);
            });

        return true;
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

app.factory('log', function($resource){
    return $resource('/api/v1/project/:project/build/:id/log/:stage',
                     {project:'@project', id: '@id', stage: '@stage'}, {
                     get: {
                         method: "GET",
                         transformResponse: function(data) {
                             return { log: data }
                         }
                     }
        });
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
    return function(val) {
        return $sce.trustAsHtml(new AnsiUp().ansi_to_html(val));
    }
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

app.filter('split', function() {
    return function (input, line) {
        if (line === undefined)
            return input.split("\n");
        else
            return input.split("\n")[line];
    };
});

app.filter('length', function() {
    return function (input) {
        return input.length;
    };

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
                '<button type="button" class="btn btn-default" style="margin-right: 8px" ng-show="show" ng-click="no()">' +
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

app.directive("ngButton", function($compile) {
    return {
        scope: {
            ngButton: '&'
        },
        link: function (scope, element, attrs) {
            // variable initialization
            scope.errorText = "";
            scope.text = element.text();
            scope.classes = element[0].className;

            // add error message popover
            attrs.$set("popover-title", "Internal Server Error");
            attrs.$set("popover-trigger", "'none'");
            attrs.$set("popover-is-open", "errorText !== ''");
            attrs.$set("popover-append-to-body", true);
            attrs.$set("uib-popover", "{{errorText}}");

            // setup eventhandlers
            attrs.$set('ng-click', 'click()');
            attrs.$set('ng-blur', 'blur()');

            // avoid endless recursion due to $scompile
            element.removeAttr("ngButton");
            element.removeAttr("ng-button");
            $compile(element)(scope);

            // local functions
            var setClass = function(className) {
                element.removeClass("btn-default");
                element.removeClass("btn-danger");
                element.removeClass("btn-success");
                element.removeClass("btn-warning");
                element.addClass(className);
            };

            // eventhandler: click()
            scope.click = function () {
                if (scope.ngButton()(scope))
                {
                    setClass("btn-warning");
                    element.html('<i class="fa fa-cog fa-spin"></i>');
                }
            };

            // eventhandler: blur()
            scope.blur = function() {
                scope.reset();
            };

            // functions for use within callback
            scope.success = function() {
                setClass("btn-success");
                element.html('<i class="fa fa-check"></i>&nbsp;Success');
            };

            scope.error = function(text) {
                setClass("btn-danger");
                element.html('<i class="fa fa-exclamation-triangle"></i>&nbsp;Error');
                scope.errorText = text;
            };

            scope.reset = function() {
                scope.errorText = "";
                element[0].className = scope.classes;
                element.html(scope.text);
            };
        }
    }
});

// ------------------------------------------------------------------------------------------------
app.service('authInterceptor', function($q) {
    var service = this;

    service.responseError = function(response) {
        if (response.status === 401){
            window.location = "/login";
        }
        return $q.reject(response);
    };
});