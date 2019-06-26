/**
 * @author Jordan
 */

// START Breadcrumb trail for webform analytics
(function () {
  // Get current page
  var currentPath = document.location.pathname;
  // Get last page from referrer; replacing host to just get path
  var lastPath = document.referrer.replace(location.origin, '');

  // If not same page as last and navigation.type was not a refresh
  if (currentPath !== lastPath && performance.navigation.type === 0) {
    var cookieName = 'breadcrumbs';
    // Use session storage to read and update cookie; because it is more consistent
    var seshCookie = sessionStorage.getItem(cookieName);

    if (seshCookie === null) {
      sessionStorage.setItem(cookieName, currentPath);
      // Update the null value to current path
      seshCookie = currentPath;
    } else {
      // If cookie not empty append current path to original value
      seshCookie += '>' + currentPath;
      sessionStorage.setItem(cookieName, seshCookie);
    }
  }

  // Populate hidden seshCookie input's if on form page
  if (document.getElementsByName('seshCookie').length > 0) {
    // Loop through all seshCookie fields, there may be multiple
    for (var i =0; i < document.getElementsByName('seshCookie').length; i++) {
      document.getElementsByName('seshCookie')[i].value = sessionStorage.getItem('breadcrumbs');
      }
  }
})();
// END Breadcrumbs