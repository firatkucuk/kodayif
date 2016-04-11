'use strict';


// #####################################################################################################################

var OPERATION_FILE_EXISTS   = {id: 1, label: 'File Exists', tag: 'fe'};
var OPERATION_FILE_CONTAINS = {id: 2, label: 'File Contains', tag: 'fc'};


// #####################################################################################################################

function kodayifController($scope, $http, $interval) {

  // ---------------------------------------------------------------------------

  $scope.getServerStatus = function() {

    var req = {
      method : 'GET',
      url    : '/status/' + $scope.uuid,
      timeout: 2000
    };

    $http(req)
      .success(function (responseBody) {

        $scope.statusItems = responseBody;
      })
      .error(function (responseBody) {
        // pass
      });
  }

  // ---------------------------------------------------------------------------

  $scope.$on('$destroy', function () {

    $interval.cancel($scope.statusFetchInterval);
    $scope.statusFetchInterval = undefined;
  });

  // ---------------------------------------------------------------------------

  $scope.send = function () {

    $scope.operationLock = true;

    if ($scope.statusFetchInterval) {
      $interval.cancel($scope.statusFetchInterval);
      $scope.statusFetchInterval = undefined;
    }

    var req = {
      method : 'POST',
      url    : '/send',
      timeout: 2000,
      data   : {
        operation : $scope.formData.operation.id,
        filePath  : $scope.formData.filePath,
        term      : $scope.formData.term
      }
    };

    $http(req)
      .success(function (responseBody) {

        $scope.uuid                = responseBody.Uuid;
        $scope.operationLock       = false;
        $scope.statusFetchInterval = $interval(function () {
          $scope.getServerStatus();
        }, 10000);

        $scope.getServerStatus();
      })
      .error(function (responseBody) {

        $scope.operationLock = false;
      });
  }

  // ---------------------------------------------------------------------------

  $scope.operationLock = false;
  $scope.formData      = {
    operation : OPERATION_FILE_EXISTS,
    filePath  : '/etc/hosts',
    term      : '8.8.8.8'
  };
  $scope.operations    = [
    OPERATION_FILE_EXISTS,
    OPERATION_FILE_CONTAINS
  ];
}


// #####################################################################################################################

(function (angular) {

  angular
    .module('kodayifApp', [])
    .controller('KodayifController', ['$scope', '$http', '$interval', kodayifController]);
})(window.angular);
