<script>
  (function() {
    // Cookies to Track Session Time
    var now = new Date();
    // var exp = new Date(now.getTime() + 7 * 24 * 3600 * 1000); // Expire in 7 days

    // Set cookie with name(key) and value
    function setCookie(name, value)
    {
      // Store as stringified JSON
      value = JSON.stringify(value);
      document.cookie=name + "=" + escape(value);
    }

    // Retrieve cookie value by name(key)
    function getCookie(name)
    {
      var reg = new RegExp(name + "=([^;]+)");
      var val = reg.exec(document.cookie);
      // Return parsed JSON if not null
      return (val != null) ? JSON.parse(unescape(val[1])) : null;
    }
   	
    function checkUserCookie(cookieName, usrCookie)
    {
      // Set cookie if one doesn't exist yet
      if (usrCookie == null) {
        var cookieVal = {
          date: now,
          pageCount: 1,
          eventFired: false
        };
        setCookie(cookieName, cookieVal);
      } else {
        // If we are at 2 or more pages and an event hasn't already been fired
        if (usrCookie.pageCount >= 2 && usrCookie.eventFired == false) {
          var cookieTime = new Date(usrCookie.date);
          var currentTime = new Date();
          var diffSeconds = (currentTime - cookieTime)/1000;

          // Then if we hit the 90sec mark we can push the event
          if (diffSeconds >= 90) {
            // Push trigger event
            window.dataLayer = window.dataLayer || [];
            window.dataLayer.push({
              event: 'siteEngagementEvent',
              attributes: {
                eventFired: '2p90s'
              }
            });

            // Set eventFired value true and update cookie
            usrCookie.eventFired = true;
            setCookie(cookieName, usrCookie);
          }
        } else if (usrCookie.pageCount < 2) {
          // Update count
          usrCookie.pageCount++;
          setCookie(cookieName, usrCookie);
        }
      }
    }
    
    var cookieName = 'siteEngagement';
    // Get user cookie
    var usrCookie = getCookie(cookieName);
    
    // Check every 5sec
    var interval = setInterval(function(){
      checkUserCookie(cookieName, usrCookie);
      
      // Stop looping if event was fired, or usrCookie is null
      if (usrCookie == null) {
        clearInterval(interval);
      }
      else if (usrCookie.eventFired == true) {
        clearInterval(interval);
      }
    }, 5000);
    
  })();
</script>