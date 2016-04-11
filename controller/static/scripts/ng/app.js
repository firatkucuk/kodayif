'use strict';


// #####################################################################################################################

var OPERATION_FILE_EXISTS   = {id: 1, label: 'File Exists', tag: 'fe'};
var OPERATION_FILE_CONTAINS = {id: 2, label: 'File Contains', tag: 'fc'};


// #####################################################################################################################

function kodayifController($scope, $http) {

  $scope.send = function () {

    $scope.operationLock = true;

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
      .success(function () {

        $scope.operationLock = false;
      })
      .error(function (responseBody) {

        $scope.operationLock = false;

        // if (responseBody && responseBody.text) {
        //   $scope.message = responseBody.text;
        // } else {
        //   $scope.message = "Couldn't communicate with device.";
        // }
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
    .controller('KodayifController', ['$scope', '$http', kodayifController]);
})(window.angular);
