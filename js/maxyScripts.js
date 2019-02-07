//checkNumber: checks and updates phone number based on page
var $ = jQuery;

function checkNumber() {
    document.getElementById('facilityNumber').innerHTML = $('.action-phone').html();
}
$(window).ready(function () {
    checkNumber();
});
//End checkNumber

//borderColor: gets color of element and applies it to whatever.
function borderColor() {
    $('#vcbMe5905489849').css("border-color", $('#facilityNumber').css("color"));
}
$(window).ready(function () {
    borderColor();
});
//End borderColor

//Override Hover background
function removeHoverEffect() {
    var originalBackground = $('#lightboxBtn').css('background');
    $('#lightboxBtn').hover(function () {
        $(this).css("background", originalBackground)
    });
}
$(window).ready(function () {
    checkNumber();
    removeHoverEffect();
});

// Accordion: checks for h2 tags and makes the p tags underneath each fold into Accordion
<style >
    /* Style the buttons that are used to open and close the accordion panel */
    button.accordion {
        background - color: #eee;
        color: #444;
      cursor: pointer;
      padding: 18px;
      width: 100%;
      text-align: left;
      border: none;
      outline: none;
      transition: 0.4s;
      margin-bottom: 4px;
      margin-top: 10px;
  }

  /* Add a background color to the button if it is clicked on (add the .active class with JS), and when you move the mouse over it (hover) */
  button.accordion.active, button.accordion:hover {
      background-color: # ddd;
    }

div.panel {
        padding: 0 18 px;
        background - color: white;
        max - height: 0;
        overflow: hidden;
        transition: max - height 0.2 s ease - out;
    } <
    /style>

var $ = jQuery;

function createAccord() {
    var x = $("article");
    var j;
    var target = "p",
        invert = ':not(' + target + ')';
    var y = x[1].getElementsByTagName("h2");
    var p = x[1].getElementsByTagName(target);
    for (j = 0; j < y.length; j++) {
        $(y[j]).nextUntil("h2").wrapAll("<div class='panel' />");
        $(y[j]).wrap("<button class=\"accordion\">");
    }
}

function activateAccord() {
    var acc = document.getElementsByClassName("accordion");
    var i;

    for (i = 0; i < acc.length; i++) {
        acc[i].onclick = function () {
            this.classList.toggle("active");
            var panel = this.nextElementSibling;
            if (panel.style.maxHeight) {
                panel.style.maxHeight = null;
            } else {
                panel.style.maxHeight = panel.scrollHeight + "px";
            }
        }
    }
}

$(document).ready(function () {
    createAccord();
    activateAccord();
});
// END Accordion

// Set Timeout for Maxy Action
var $ = jQuery;

$(document).ready(function () {
    setTimeout(function () {
        actions
            .set('TimerAction', '1')
            .send()
    }, 90000);
});
// End Set Timeout

// Cookies to Track Session Time
var now = new Date();
var exp = new Date(now.getTime() + 7 * 24 * 3600 * 1000); // Expire in 7 days

function setCookie(name, value) // Set cookie with name(key) and value
{
    document.cookie = name + "=" + escape(value) + "; path=/; expires=" + exp.toString();
}

function getCookie(name) // Retrieve cookie value by name(key)
{
    var reg = new RegExp(name + "=([^;]+)");
    var val = reg.exec(document.cookie);
    return (val != null) ? unescape(val[1]) : null;
}
// END Cookies Track Session Time

// Check if variant lightbox or whatever is on screen before attributing call actions
// This script checks if the variant callout/lightbox id is visible on the page.
// If it is then it will attribute call actions if a call is made.
// Must add "CallsVariantExposure" to campaigns to see actions.
var $ = jQuery;
var checkVar = false;

$(window).scroll(function (event) {
    if ($(".vcb-me-content-wrapper").is(':visible') && checkVar == false) {
        window.MaxymiserDidRun = false;
        window.MaxymiserFailed = false;
        var call = document.getElementById("call");
        if (call) {
            actions.send('CallsVariantExposure', call.getAttribute("data-value"))
                .done(function () {
                    window.MaxymiserDidRun = true;
                })
                .fail(function () {
                    window.MaxymiserFailed = true;
                });
        } else {
            window.MaxymiserFailed = true;
        }
        checkVar = true;
        console.log("Call Action Set.");
    }
});
// END Call Action Check

// START Video Plays Action
if (window.location.host != "www.life-healing.com" && window.location.host != "www.wellnessresourcecenter.com") {
    when(function () {
        return window.jQuery;
    }).done(function () {
        var $ = jQuery;

        $('.acadia_video_inner iframe').load(function () {

            if (window.location.pathname.indexOf("/lp/ar1/") !== -1) {
                // Don't fire on landing page content tests!
                return;
            }
            // grab page type to pass as action attribute
            var pageType = '';

            // track engagement with video
            var player = $('.acadia_video_inner iframe');
            var playerOrigin = '*';
            //window.addEventListener('message', onMessageReceived, false);

            // Helper function for sending a message to the player
            function post(action, value) {
                var data = {
                    method: action
                };
                if (value) {
                    data.value = value;
                }
                var message = JSON.stringify(data);
                player[0].contentWindow.postMessage(message, playerOrigin);
            }
            // handle messages received from the player
            function onMessageReceived(event) {
                console.log(event);
                // handle messages from the vimeo player only
                if (!(/^https?:\/\/player.vimeo.com/).test(event.origin)) {
                    return false;
                }
                if (playerOrigin == '*') {
                    playerOrigin = event.origin;
                }
                var videoData = JSON.parse(event.data);

                post('addEventListener', 'play');
                post('addEventListener', 'pause');
                post('addEventListener', 'finish');

                if (videoData.event == 'play') {
                    actions.send('VideoPlay', 1, pageType);
                }
                if (videoData.event == 'pause') {
                    var timePlayed = Math.round(videoData.data.seconds);
                    actions.send('VideoPause', timePlayed, pageType);
                }
                if (videoData.event == 'finish') {
                    actions.send('VideoFinish', 1, pageType);
                }
            }
            // listen for messages from the player
            if (window.addEventListener) {
                window.addEventListener('message', onMessageReceived, false);
            } else {
                window.attachEvent('onmessage', onMessageReceived, false);
            }
        });
    });
}
// END Video plays action

// START Breadcrumb trail for webform analytics
(function() {
    // Set cookie with name(key) and value
    function setCookie(name, value) {
      document.cookie = name + "=" + value;
    }
  
    // Retrieve cookie value by name(key)
    function getCookie(name) {
      var reg = new RegExp(name + "=([^;]+)");
      var val = reg.exec(document.cookie);
      // Return parsed JSON if not null
      return (val != null) ? val[1] : null;
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