'use strict';

angular.module('myApp.view6', ['ngRoute', 'myApp'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view6', {
    templateUrl: 'view6/view6.html',
    controller: 'View6Ctrl'
  });
}])

.controller('View6Ctrl', function($scope, productService, $http, $location) {
	$('body').removeClass('modal-open');
	$('.modal-backdrop').remove();
	$scope.username = productService.getProducts()[0];
	//$scope.username = "cyrus";

	
});