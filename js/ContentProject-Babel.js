'use strict';

function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
        throw new TypeError("Cannot call a class as a function");
    }
}

var Person = function () {

    // URLSearchParams replacement func for stupid IE.
    Person.prototype.getParameterByName = function getParameterByName(name, url) {
        if (!url) url = window.location.href;
        name = name.replace(/[\[\]]/g, "\\$&");
        var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
            results = regex.exec(url);
        if (!results) return null;
        if (!results[2]) return '';
        return decodeURIComponent(results[2].replace(/\+/g, " "));
    }

    // Constructor
    function Person() {
        _classCallCheck(this, Person);

        // Get query string object for parsing
        //this.qString = new URLSearchParams(window.location.search);
        //this.gender = this.getParameterByName('utm_content').toLowerCase(); // Get gender from query string
        this.drug = window.location.search.toLowerCase(); // Set drug to whole query string so we can search the whole thing in switch
        this.gender = this.drug.indexOf('female') !== -1 ? "female" : this.drug.indexOf('male') !== -1 ? "male" : "default";
        //this.term = this.getParameterByName('utm_term').toLowerCase(); // Get term for drug
        this.campaign = this.getParameterByName('utm_campaign').toLowerCase(); // Get campaign for any extra content variables
        this.facility = document.getElementById('utm_facility').innerHTML; // Retrieve Facility from maxy tag
        this.h1 = "";
        this.p = "";
        this.vidReplace = "";

        // Look for Gen 2 video first
        if (document.getElementsByClassName('acadia_video_inner')[0]) { 
            this.video = document.getElementsByClassName('acadia_video_inner')[0];
        } else {
            // Otherwise get video by Gen 1 class
            this.video = document.getElementsByClassName('video')[0];
        }
        
    }

    // Replace <h1> and <p> content with generated campaign

    Person.prototype.updateContent = function updateContent() {
        try {
            // Use new header content input to overwrite default h1, "utm_h1" id inserted from maxymiser
            document.getElementById("utm_h1").innerHTML = this.h1;
            // Use new p content input to overwrite default p, "utm_p" id inserted from maxymiser
            document.getElementById("utm_p").innerHTML = this.p;
            // Update video content
            this.video.innerHTML = this.vidReplace; 
        } catch (err) {
            console.log(err.message);
        }
    };

    // Capitalize func for drug variables

    Person.prototype.capitalize = function capitalize(s) {
        return s.charAt(0).toUpperCase() + s.slice(1);
    };

    // START All Drug Contents

    Person.prototype.alcoholContent = function alcoholContent() {
        // Female Alcohol
        if (this.gender == "female") {
            this.h1 = "Alcoholism shouldn’t be a family disease";
            this.p = "Alcohol abuse puts financial and physical stress on everyone. " + this.facility + " understands recovery and will commit to you and your loved one. Call and get started.";
        }
        // Male Alcohol
        else if (this.gender == "male") {
            this.h1 = "This is Your Last Stop Before Rock Bottom ";
            this.p = "Alcoholism loses its power when it stops being a secret. " + "Stay at " + this.facility + " and join a community committed to your recovery. " + "Build the foundation you need to achieve long-term success and unravel the underlying cause of your abuse. " + "Addiction is not a choice, or a matter of weak will—with our help, you can succeed, just call.";
        }
        // Default Alcohol
        else {
            this.h1 = "Alcoholism Knocked You Down. We'll Help You Back Up.";
            this.p = "Alcohol abuse is not easy to overcome. " + this.facility + " will meet you where you are and help you recover. " + "Our alcohol addiction treatment program focuses on unraveling the underlying causes of addiction. " + "Learn to trust yourself again. Get started and call us today.";
        }
    }; // END Alcohol Content

    Person.prototype.opioidContent = function opioidContent(opioid) {
        // Female Opioid
        if (this.gender == "female") {
            this.h1 = this.capitalize(opioid) + " Addiction Is Hard To Beat On Your Own";
            this.p = this.capitalize(opioid) + " addiction compromises your relationships. Isolation makes it worse. " + "Join our recovery community - we can help. " + this.facility + " specializes in extinguishing the triggers that cause " + opioid + " abuse. Call today to set up an appointment.";
        }
        // Male Opioid
        else if (this.gender == "male") {
            this.h1 = "Conquer Your " + this.capitalize(opioid) + " Addiction.";
            this.p = this.capitalize(opioid) + " addiction is a secret that puts you in isolation. Set yourself free and defeat your demons. " + this.facility + " specializes in unraveling the underlying causes of your addiction. Call and get help.";
        }
        // Default Opioid
        else {
            this.h1 = "Casual " + this.capitalize(opioid) + " Use Doesn’t Exist";
            this.p = "When " + opioid + " use is present, addiction follows. " + this.facility + " can help you get clean and stay that way. " + "Take your life back from " + opioid + " addiction. Call and get help.";
        }
    }; // END Opioid Content

    Person.prototype.painkillerContent = function painkillerContent(painkiller) {
        // Female Painkiller
        if (this.gender == "female") {
            this.h1 = "Accidental Addiction is Still Addiction";
            this.p = this.capitalize(painkiller) + " addiction infects everyone it touches. Regardless of how it starts, " + "the sooner you treat it the closer you are to regaining control of your life. " + this.facility + " provides treatment required to achieve lasting recovery. Let us help you. Call and get started.";
        }
        // Male Painkiller
        else if (this.gender == "male") {
            this.h1 = "Manage Physical Pain Without " + this.capitalize(painkiller);
            this.p = "Make the choice to build a strong foundation for lasting recovery. " + this.facility + " is here to help you conquer " + painkiller + " while managing your pain. Get your life back. Call Today.";
        }
        // Default Painkiller
        else {
            this.h1 = "Manage Your Pain Without Addiction.";
            this.p = "Dependency on " + painkiller + " is a hard thing to kick, but it's worth the struggle. " + this.facility + " is here to help you get back on your feet. " + "Getting hurt wasn’t your fault, but getting better is your responsibility. Call and get started.";
        }
    }; // END Painkiller Content

    Person.prototype.stimulantContent = function stimulantContent(stimulant) {
        // Female Stimulant
        if (this.gender == "female") {
            this.h1 = "Recovery Begins The Sooner You Act";
            this.p = "Your loved one needs you. Tough love is still love. Don't let fear prevent you from saving their life. " + this.facility + " offers comprehensive " + stimulant + " addiction treatment and our team will be with you every step of the way. " + "Call and get started.";
        }
        // Male Stimulant
        else if (this.gender == "male") {
            this.h1 = "Your Next Choice Could Save Your Life.";
            this.p = "You don't need " + stimulant + " to make you feel invincible. At " + this.facility + " we help you build a strong foundation for lasting recovery. " + "Come meet our experienced team who will commit to your treatment. Learn to recognize your addiction triggers and adapt. " + "You can evolve past " + stimulant + ". Call us today.";
        }
        // Default Stimulant
        else {
            this.h1 = "Start Recovery. Get Your Life Back From " + this.capitalize(stimulant);
            this.p = "You’re tougher than you think — conquer " + stimulant + ". Stay at " + this.facility + " and meet people who will commit to your success. You can build a strong foundation for recovery. " + "Learn to recognize your triggers and grow from your experience. Call and get started today.";
        }
    }; // END Stimulant Content

    Person.prototype.sedativeContent = function sedativeContent(sedative) {
        // Female Sedative
        if (this.gender == "female") {
            this.h1 = "Stop the Numbing. Start Recovery.";
            this.p = this.capitalize(sedative) + " abuse puts financial and emotional stress on everyone. " + this.facility + " understands recovery and will commit to you and your loved one. " + "The first step towards recovery is just a phone call away.";
        }
        // Male Sedative
        else if (this.gender == "male") {
            this.h1 = "Numbness Is Temporary. Recovery Can Last.";
            this.p = this.capitalize(sedative) + " abuse may seem like a fix. But after you wake up, your problems are still there. " + this.facility + " will work with you to address those problems head-on, opening the door to long-term recovery. " + "You can succeed with our help. Just call.";
        }
        // Default Sedative
        else {
            this.h1 = "Hiding " + this.capitalize(sedative) + " Abuse Only Hides The Problem.";
            this.p = "Are you ready to start truly living again? " + this.facility + " will meet you where you are and help you recover. " + "Our " + sedative + " addiction treatment focuses on unraveling the underlying causes of addiction. " + "You can learn to trust yourself again. Get started and call us today.";
        }
    }; // END Sedative Content

    Person.prototype.generalDrugContent = function generalDrugContent() {
        // Female General Drug
        if (this.gender == "female") {
            this.h1 = "Recovery Begins The Sooner You Act";
            this.p = "Your loved one needs you. Tough love is still love. Don't let fear prevent you from saving their life. " + this.facility + " offers comprehensive addiction treatment and our team will be with you every step of the way. " + "Call and get started.";
        }
        // Male General Drug
        else if (this.gender == "male") {
            this.h1 = "Your Next Choice Could Save Your Life.";
            this.p = "You don't need drugs to make you feel invincible. At " + this.facility + " we help you build a strong foundation for lasting " + "recovery. Come meet our experienced team who will commit to your treatment. " + "Learn to recognize your addiction triggers and adapt. You can evolve past addiction. Call us today.";
        }
        // Default General Drug
        else {
            this.h1 = "Start Recovery. Get Your Life Back From Drugs.";
            this.p = "You’re tougher than you think—conquer addiction. Stay at " + this.facility + " and meet people who will commit to your success. " + "You can build a strong foundation for recovery. Learn to recognize your triggers and grow from your experience. " + "Call and get started today.";
        }
    }; // END General Drug Content
    // END All Drug Contents

    // Massive switch to fall through 'utm_content' value and generate

    Person.prototype.parseDrug = function parseDrug() {
        switch (true) {
            // START Alcohol
            case this.drug.indexOf("alcohol") !== -1:
                this.alcoholContent();
                 // Update Video Content
                if (this.gender == "female") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241237917" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else if (this.gender == "male") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241182918" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/240844156" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                }
                break;
                // END Alcohol

                // START Opiods
            case this.drug.indexOf("heroin") !== -1:
                this.opioidContent("heroin");
                if (this.gender == "female") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241233680" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else if (this.gender == "male") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241200739" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/240844000" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                }
                break;
            case this.drug.indexOf("opiate") !== -1:
                this.opioidContent("opiate");
                break;
            case this.drug.indexOf("opioid") !== -1:
                this.opioidContent("opioid");
                break;
                // END Opiods

                // START Stimulants
            case this.drug.indexOf("meth") !== -1:
                this.stimulantContent("meth");
                break;
            case this.drug.indexOf("amphetamine") !== -1:
                this.stimulantContent("amphetamine");
                break;
            case this.drug.indexOf("pcp") !== -1:
                this.stimulantContent("PCP");
                break;
            case this.drug.indexOf("stimulant") !== -1:
                this.stimulantContent("stimulants");
                break;
                // END Stimulants

                // START Cocaine
            case this.drug.indexOf("cocaine") !== -1:
                this.stimulantContent("cocaine");
                if (this.gender == "female") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241236618" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else if (this.gender == "male") {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/241195195" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                } else {
                    this.vidReplace = '<iframe src="https://player.vimeo.com/video/240844853" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
                }
                break;
                // END Cocaine

                // START Sedatives
            case this.drug.indexOf("xanax") !== -1:
                this.sedativeContent("Xanax");
                break;
            case this.drug.indexOf("klonopin") !== -1:
                this.sedativeContent("Klonopin");
                break;
            case this.drug.indexOf("barbiturate") !== -1:
                this.sedativeContent("barbiturate");
                break;
            case this.drug.indexOf("ativan") !== -1:
                this.sedativeContent("Ativan");
                break;
            case this.drug.indexOf("benzo") !== -1:
                this.sedativeContent("benzo");
                break;
            case this.drug.indexOf("valium") !== -1:
                this.sedativeContent("Valium");
                break;
            case this.drug.indexOf("sedative") !== -1:
                this.sedativeContent("sedative");
                break;
                // END Sedatives

                // START Painkillers
            case this.drug.indexOf("percocet") !== -1:
                this.painkillerContent("Percocet");
                break;
            case this.drug.indexOf("codeine") !== -1:
                this.painkillerContent("codeine");
                break;
            case this.drug.indexOf("roxicodone") !== -1:
                this.painkillerContent("Roxicodone");
                break;
            case this.drug.indexOf("hydrocodone") !== -1:
                this.painkillerContent("hydrocodone");
                break;
            case this.drug.indexOf("oxycodone") !== -1:
                this.painkillerContent("oxycodone");
                break;
            case this.drug.indexOf("oxycontin") !== -1:
                this.painkillerContent("OxyContin");
                break;
            case this.drug.indexOf("oxy") !== -1:
                this.painkillerContent("oxy");
                break;
            case this.drug.indexOf("vicodin") !== -1:
                this.painkillerContent("Vicodin");
                break;
            case this.drug.indexOf("demerol") !== -1:
                this.painkillerContent("Demerol");
                break;
            case this.drug.indexOf("morphine") !== -1:
                this.painkillerContent("morphine");
                break;
            case this.drug.indexOf("fentanyl") !== -1:
                this.painkillerContent("fentanyl");
                break;
            case this.drug.indexOf("lortab") !== -1:
                this.painkillerContent("Lortab");
                break;
            case this.drug.indexOf("dilaudid") !== -1:
                this.painkillerContent("Dilaudid");
                break;
            case this.drug.indexOf("painkiller") !== -1:
                this.painkillerContent("painkillers");
                break;
                // END Painkillers

                // Default is "Generic Drug"
            default:
                this.generalDrugContent();
                break;
        }

        // If not Alcohol, Heroin, or Cocaine, video replacement will be empty. Fill it with generic video
        if (this.vidReplace.length < 1) {
            if (this.gender == "female") {
                this.vidReplace = '<iframe src="https://player.vimeo.com/video/241225888" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
            } else if (this.gender == "male") {
                this.vidReplace = '<iframe src="https://player.vimeo.com/video/241212117" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
            } else {
                this.vidReplace = '<iframe src="https://player.vimeo.com/video/240205224" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';
            }
        }

        // After switch logic sets appropriate h1 and p, call func to edit document
        return this.updateContent();
    }; // END parseDrug()


    // Callable to generate content after class construction per facility
    Person.prototype.generateCampaign = function generateCampaign() {
        // Find drug content to display
        this.parseDrug();
    };

    return Person;
}(); // END Person Class


// Check for jQuery
function checkJquery(count) {
    if (window.jQuery) {
        var $ = jQuery;
        console.log("jQuery found");
        var $ = jQuery;
        // Content Generate
        $(document).ready(function () {
            console.log("utm check fire");
            var h1ElementExists = document.getElementById("utm_h1");
            var pElementExists = document.getElementById("utm_p");
            if (h1ElementExists && pElementExists) {
                // If we can find the "utm" tags, it's a valid campaign so generate content
                var person = new Person();
                person.generateCampaign();
                console.log("utm content loaded");
            }
        });

    } else if (count >= 20) {
        var script = document.createElement("SCRIPT");
        script.src = 'https://ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js';
        script.type = 'text/javascript';
        script.onload = function () {
            var $ = window.jQuery;
            var $ = jQuery;
            // Content Generate
            $(document).ready(function () {
                console.log("utm check fire");
                var h1ElementExists = document.getElementById("utm_h1");
                var pElementExists = document.getElementById("utm_p");
                if (h1ElementExists && pElementExists) {
                    // If we can find the "utm" tags, it's a valid campaign so generate content
                    var person = new Person();
                    person.generateCampaign();
                    console.log("utm content loaded");
                }
            });
        };
        document.getElementsByTagName("head")[0].appendChild(script);
    } else {
        count += 1;
        console.log("count = "+count);
        setTimeout(function () {checkJquery(count)}, 200);

    }
}

try {
    var $ = jQuery;
    // Generate content on load
    $(document).ready(function () {
        setTimeout(function () {
            // START timeout to wait for maxy tags to exist
            console.log("utm check fire");
            var h1ElementExists = document.getElementById("utm_h1");
            var pElementExists = document.getElementById("utm_p");
            if (h1ElementExists && pElementExists) {
                // If we can find the "utm" tags, it's a valid campaign so generate content
                var person = new Person();
                person.generateCampaign();
                console.log("utm content loaded");
            }

            // END timeout
        }, 10);
    });
} catch (err) {
    checkJquery(0);
    console.log('timeout');
}