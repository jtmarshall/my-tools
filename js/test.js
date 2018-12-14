function runZopim(timer) {
    if (typeof $zopim !== 'undefined') {
        $zopim(function () {
            if (zendeskHide == 1) {
                $zopim.livechat.hideAll();
            } else {
                if (zendeskReposition == 1) { // zendeskReposition is set in /inc/theme-dynamic-css.php
                    $zopim.livechat.button.setOffsetVertical(50);
                    $zopim.livechat.button.setOffsetVerticalMobile(50);
                    $zopim.livechat.window.setOffsetVertical(50);
                }
                $zopim.livechat.theme.setTheme('simple'); // this swaps the questions bubble with a smaller panel
                $zopim.livechat.theme.setColor(zendeskPrimaryColor); // zendeskPrimaryColor is set in /inc/theme-dynamic-css.php
                // $zopim.livechat.bubble.setColor(zendeskSecondaryColor); // zendeskSecondaryColor is set in /inc/theme-dynamic-css.php
                $zopim.livechat.badge.setColor(zendeskSecondaryColor); // zendeskSecondaryColor is set in /inc/theme-dynamic-css.php
                $zopim.livechat.badge.setLayout('text_only');
                $zopim.livechat.badge.setText('Have questions? Click to chat.');
                $zopim.livechat.theme.setFontConfig({
                    google: {
                        families: [zendeskFont]
                    }
                }, zendeskFont); // zendeskFont is set in /inc/theme-dynamic-css.php
                $zopim.livechat.theme.reload();
                // $zopim.livechat.badge.show();
            }
        });
    } else if (timer < 20) {
        setTimeout(function () {
            // Increment
            timer += 1;
            // Wait 100ms for $zopim to exist, try again
            runZopim(timer);
        }, 100);
    } else {
        // Return after 20 failed intervals
        console.log("zopim not loaded");
        return;
    }
}

runZopim(0);


// try {
//     var $ = jQuery;
//     // Generate content on load
//     $(document).ready(function () {
//         setTimeout(function () {
//             // START timeout to wait for maxy tags to exist
//             console.log("utm check fire");
//             var h1ElementExists = document.getElementById("utm_h1");
//             var pElementExists = document.getElementById("utm_p");
//             if (h1ElementExists && pElementExists) {
//                 // If we can find the "utm" tags, it's a valid campaign so generate content
//                 var person = new Person();
//                 person.generateCampaign();
//                 console.log("utm content loaded");
//             }

//             // END timeout
//         }, 10);
//     });
// } catch (err) {
//     runZopim();
// }