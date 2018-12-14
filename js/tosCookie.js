/* 
AUTHOR: Jordan Marshall
Small JS lib that can be added to sites and used to trigger t minute Time On Site (TOS) mico-conversions.
Uses AcadiaTOS singleton class to get/set cookie, and to load conversion func array
*/

var AcadiaTOS = (function () {
    var instance; // prevent modification of "instance" variable
    var reg_funcs = {}; // Registered Functions dictionary
    var tos = Math.abs(new Date() - Date.parse(getTOSCookie())); // Difference between TOS cookie and now
    var seconds = Math.round((tos / 1000)); // Convert diff in time to seconds

    // Set AcadiaTOS cookie
    function setTOSCookie() {
        var time = new Date();
        var domain = document.domain;
        document.cookie = "AcadiaTOS" + "=" + time + ";domain=" + domain + ";path=/";
    }

    // Looks for the AcadiaTOS cookie and returns it if exists
    function getTOSCookie() {
        var name = "AcadiaTOS" + "=";
        var decodedCookie = decodeURIComponent(document.cookie);
        var ca = decodedCookie.split(';');
        for (var i = 0; i < ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0) == ' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) == 0) {
                return c.substring(name.length, c.length);
            }
        }
        return "";
    }

    // New instance
    function createInstance() {
        setTOSCookie(); // Create TOS cookie for instance
        reg_funcs = {};
        return reg_funcs;
    }

    // Func calls for 'AcadiaTOS' singleton
    return {
        getInstance: function () {
            if (!instance) {
                instance = createInstance();
            }
            return instance;
        },
        getRegisteredFuncs: function () {
            return reg_funcs;
        },
        getTOS: function () {
            if (getTOSCookie().length < 1) {
                setTOSCookie();
            } else {
                tos = Math.abs(new Date() - Date.parse(getTOSCookie())); // Difference between TOS cookie and now
                seconds = Math.round((tos / 1000));
                return seconds;
            }
        },
        addConversion: function (time, callable) {
            reg_funcs[time] = callable;
        },
        removeConversion: function (index) {
            delete reg_funcs[index];
            console.log(reg_funcs[index]);
        }
    };
})();


// To Run the singleton
function runAcadiaTOS() {
    // var foo = function () {
    //     console.log("test");
    // }
    
    // console.log(AcadiaTOS.getTOS());

    // AcadiaTOS.addConversion(10, foo);
    // AcadiaTOS.addConversion(30, foo);
    // AcadiaTOS.addConversion(80, foo);
    // console.log(AcadiaTOS.getRegisteredFuncs());

    //AcadiaTOS.getInstance();

    // Retrieve our Class variables for this instance
    var seconds = AcadiaTOS.getTOS();
    var funcs = AcadiaTOS.getRegisteredFuncs();
    var upperLimit = 400;

    // Loop func to substitute 'while' so we can use timeout function
    var looper = function () {
        if (seconds < upperLimit) {
            setTimeout(function () { // call 1s setTimeout when loop is called
                // Re-Get seconds and registered funcs
                seconds = AcadiaTOS.getTOS();
                funcs = AcadiaTOS.getRegisteredFuncs();

                // Check if there are registered funcs, otherwise skip iteration
                if (Object.keys(funcs).length > 0) {
                    // Iterate through registered functions to check against TOS
                    for (var index in funcs) {
                        //console.log(seconds + " - funcs length: " + Object.keys(funcs).length);

                        // Make sure our limit will always include highest time index value
                        if (upperLimit < parseInt(index)) {
                            // Tack on 5sec to allow some buffer
                            upperLimit = parseInt(index) + 10;
                        }
                        // If we get to the time of an index: call func
                        if (seconds >= parseInt(index)) {
                            console.log(seconds, index);
                            // Call func at index/time
                            funcs[index]();
                            console.log("5min TOS");
                            // Remove from registered functions array after fire
                            AcadiaTOS.removeConversion(index);
                        }
                    }
                }
                // Loop
                looper();
            }, 1000);
        } else {
            // If we reach our upperLimit break out of loop
            console.log("AcadiaTOS Exit.")
            return;
        }
    }

    // START Loop
    looper();
}

runAcadiaTOS();