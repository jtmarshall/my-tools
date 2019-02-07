/**
 * @author Jordan
 */

// START Breadcrumb trail for webform analytics
(() => {
  // Set cookie with name(key) and value
  setCookie = (name, value) => {
    document.cookie = name + "=" + value;
  }

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
    } else {
      // If cookie not empty append current path to original value
      seshCookie += '>' + currentPath;
      console.log('update cookie: ', seshCookie);
      sessionStorage.setItem(cookieName, seshCookie);
    }

    // Overwrite cookie with session storage val
    setCookie(cookieName, seshCookie);
  }
})();
// END Breadcrumbs