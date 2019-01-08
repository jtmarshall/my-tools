// Class to process query string to generate corresponding content for campaign

class Person {
    constructor() {
        // Get query string object for parsing
        this.qString = new URLSearchParams(window.location.search)
        this.gender = this.qString.get('utm_content').toLowerCase() // Get gender from query string
        this.term = this.qString.get('utm_term').toLowerCase() // Get term for drug
        this.campaign = this.qString.get('utm_campaign').toLowerCase() // Get campaign for any extra content variables
        this.drug = window.location.search.toLowerCase() // Set drug to whole query string so we can search the whole thing in switch
        this.facility = document.getElementById('utm_facility').innerHTML // Retrieve Facility from maxy tag
        this.h1 = "";
        this.p = "";
    }


    // Replace <h1> and <p> content with generated campaign
    updateContent() {
        try {
            // Use new header content input to overwrite default h1, "utm_h1" id inserted from maxymiser
            document.getElementById("utm_h1").innerHTML = this.h1;
            // Use new p content input to overwrite default p, "utm_p" id inserted from maxymiser
            document.getElementById("utm_p").innerHTML = this.p;
        }
        catch(err) {
            console.log(err.message);
        }
    }


    // Capitalize func for drug variables
    capitalize(s) {
        return s.charAt(0).toUpperCase() + s.slice(1);
    }


    // START All Drug Contents
    alcoholContent() {
        // Female Alcohol
        if (this.gender == "female") {
            this.h1 = "Alcoholism shouldn’t be a family disease";
            this.p = "Alcohol abuse puts financial and physical stress on everyone. " +
                this.facility + " understands recovery and will commit to you and your loved one. Call and get started.";
        }
        // Male Alcohol
        else if (this.gender == "male") {
            this.h1 = "This is Your Last Stop Before Rock Bottom ";
            this.p = "Alcoholism loses its power when it stops being a secret. " +
                "Stay at " + this.facility + " and join a community committed to your recovery. " +
                "Build the foundation you need to achieve long-term success and unravel the underlying cause of your abuse. " +
                "Addiction is not a choice, or a matter of weak will—with our help, you can succeed, just call.";
        }
        // Default Alcohol
        else {
            this.h1 = "Alcoholism Knocked You Down. We'll Help You Back Up.";
            this.p = "Alcohol abuse is not easy to overcome. " + this.facility + " will meet you where you are and help you recover. " +
                "Our alcohol addiction treatment program focuses on unraveling the underlying causes of addiction. " +
                "Learn to trust yourself, again, it starts with a phone call.";
        }
    } // END Alcohol Content


    opioidContent(opioid) {
        // Female Opioid
        if (this.gender == "female") {
            this.h1 = this.capitalize(opioid) + " Addiction Is Hard To Beat On Your Own";
            this.p = this.capitalize(opioid) + " addiction compromises your relationships. Isolation makes it worse. " +
                "Join our recovery community - we can help. " + this.facility +
                " specializes in extinguishing the triggers that cause " + opioid + " abuse. Call today to set up an appointment.";
        }
        // Male Opioid
        else if (this.gender == "male") {
            this.h1 = "Conquer Your " + this.capitalize(opioid) + " Addiction.";
            this.p = this.capitalize(opioid) + " addiction is a secret that puts you in isolation. Set yourself free and defeat your demons. " +
                this.facility + " specializes in unraveling the underlying causes of your addiction. Call and get help.";
        }
        // Default Opioid
        else {
            this.h1 = "Casual " + this.capitalize(opioid) + " Use Doesn’t Exist";
            this.p = "When " + opioid + " use is present, addiction follows. " + this.facility + " can help you get clean and stay that way. " +
                "Take your life back from opiate addiction. Call and get help.";
        }
    } // END Opioid Content


    painkillerContent(painkiller) {
        // Female Painkiller
        if (this.gender == "female") {
            this.h1 = "Accidental Addiction is Still Addiction";
            this.p = this.capitalize(painkiller) + " addiction infects everyone it touches. Regardless of how it starts, " +
                "the sooner you treat it the closer you are to regaining control of your life. " +
                this.facility + " provides treatment required to achieve lasting recovery. Let us help you. Call and get started.";
        }
        // Male Painkiller
        else if (this.gender == "male") {
            this.h1 = "Manage Physical Pain Without " + this.capitalize(painkiller);
            this.p = "Make the choice to build a strong foundation for lasting recovery. " +
                this.facility + " is here to help you conquer " + painkiller + " while managing your pain. Get your life back. Call Today.";
        }
        // Default Painkiller
        else {
            this.h1 = "Manage Your Pain Without Addiction.";
            this.p = "Dependency on " + painkiller + " is a hard thing to kick, but it's worth the struggle. " +
                this.facility + " is here to help you get back on your feet. " +
                "Getting hurt wasn’t your fault, but getting better is your responsibility. Call and get started.";
        }
    } // END Painkiller Content


    stimulantContent(stimulant) {
        // Female Stimulant
        if (this.gender == "female") {
            this.h1 = "Recovery Begins The Sooner You Act";
            this.p = "Your loved one needs you. Tough love is still love. Don't let fear prevent you from saving their life. " +
                this.facility + " offers comprehensive " + stimulant + " addiction treatment and our team will be with you every step of the way. " +
                "Call and get started.";
        }
        // Male Stimulant
        else if (this.gender == "male") {
            this.h1 = "Your Next Choice Could Save Your Life.";
            this.p = "You don't need " + stimulant + " to make you feel invincible. At " +
                this.facility + " we help you build a strong foundation for lasting recovery. " +
                "Come meet our experienced team who will commit to your treatment. Learn to recognize your addiction triggers and adapt. " +
                "You can evolve past " + stimulant + ". Call us today.";
        }
        // Default Stimulant
        else {
            this.h1 = "Start Recovery. Get Your Life Back From " + this.capitalize(stimulant);
            this.p = "You’re tougher than you think — conquer " + stimulant + ". Stay at " +
                this.facility + " and meet people who will commit to your success. You can build a strong foundation for recovery. " +
                "Learn to recognize your triggers and grow from your experience. Call and get started today.";
        }
    } // END Stimulant Content


    sedativeContent(sedative) {
        // Female Sedative
        if (this.gender == "female") {
            this.h1 = "Stop the Numbing. Start Recovery.";
            this.p = this.capitalize(sedative) + " abuse puts financial and emotional stress on everyone. " +
                this.facility + " understands recovery and will commit to you and your loved one. " +
                "The first step towards recovery is just a phone call away.";
        }
        // Male Sedative
        else if (this.gender == "male") {
            this.h1 = "Numbness Is Temporary. Recovery Can Last.";
            this.p = this.capitalize(sedative) + " abuse may seem like a fix. But after you wake up, your problems are still there. " +
                this.facility + " will work with you to address those problems head-on, opening the door to long-term recovery. " +
                "You can succeed with our help. Just call.";
        }
        // Default Sedative
        else {
            this.h1 = "Hiding " + this.capitalize(sedative) + " Abuse Only Hides The Problem.";
            this.p = "Are you ready to start truly living again? " + this.facility + " will meet you where you are and help you recover. " +
                "Our " + sedative + " addiction treatment focuses on unraveling the underlying causes of addiction. " +
                "You can learn to trust yourself again. Get started and call us today.";
        }
    } // END Sedative Content


    generalDrugContent() {
        // Female General Drug
        if (this.gender == "female") {
            this.h1 = "Recovery Begins The Sooner You Act";
            this.p = "Your loved one needs you. Tough love is still love. Don't let fear prevent you from saving their life. " +
                this.facility + " offers comprehensive addiction treatment and our team will be with you every step of the way. " +
                "Call and get started.";
        }
        // Male General Drug
        else if (this.gender == "male") {
            this.h1 = "Your Next Choice Could Save Your Life.";
            this.p = "You don't need drugs to make you feel invincible. At " + this.facility + " we help you build a strong foundation for lasting " +
                "recovery. Come meet our experienced team who will commit to your treatment. " +
                "Learn to recognize your addiction triggers and adapt. You can evolve past addiction. Call us today.";
        }
        // Default General Drug
        else {
            this.h1 = "Start Recovery. Get Your Life Back From Drugs.";
            this.p = "You’re tougher than you think—conquer addiction. Stay at " + this.facility + " and meet people who will commit to your success. " +
                "You can build a strong foundation for recovery. Learn to recognize your triggers and grow from your experience. " +
                "Call and get started today.";
        }
    } // END General Drug Content
    // END All Drug Contents


    // Massive switch to fall through 'utm_content' value and generate 
    parseDrug() {
        switch (true) {
            // START Alcohol
            case this.drug.indexOf("alcohol") !== -1:
                this.alcoholContent();
                break;
                // END Alcohol

                // START Opiods
            case this.drug.indexOf("heroin") !== -1:
                this.opioidContent("heroin");
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
        // After switch logic sets appropriate h1 and p, call func to edit document
        console.log("h1: "+this.h1);
        return this.updateContent();
    } // END parseDrug()


    // Callable to generate content after class construction per facility
    generateCampaign() {
        // If gender is neither male nor female, check if can find in campaign string
        if (this.gender != "female" && this.gender != "male") {
            // Banner ads will have gender in campaign section instead, so check for it
            // If female exists in string set to female, otherwise check for male, otherwise set to default
            this.gender = this.campaign.indexOf('female') !== -1 ? "female" : this.campaign.indexOf('male') !== -1 ? "male" : "default";
        }

        // Find drug content to display
        this.parseDrug();
    }

} // END Person Class


// Generate content on load
var $ = jQuery;
$(document).ready(function() {
  setTimeout(function() {
  // START timeout to wait for maxy tags to exist

    console.log("utm check fire");
    var h1ElementExists = document.getElementById("utm_h1");
    var pElementExists = document.getElementById("utm_p");
    if (h1ElementExists && pElementExists) {
        // If we can find the "utm" tags, it's a valid campaign so generate content
        var person = new Person();
        person.generateCampaign();
        console.log("utm content fire");
    }
  
  // END timeout
  }, 5);
});


?utm_source=google&utm_medium=cpc&sf_shortname=nonbrandfcrc&utm_campaign=Wilderness+T2&utm_term=cocaine%20recovery&k_clickid=40eb2efb-f7f9-4389-b980-33c2e6ab7480&utm_content=female&gclid=CI-RwZ_p7dQCFUa4wAodQ1MJUg

?utm_source=google&utm_medium=cpc&sf_shortname=nonbrandfcrc&utm_campaign=Wilderness+T2&utm_term=wilderness%20recovery&k_clickid=40eb2efb-f7f9-4389-b980-33c2e6ab7480&utm_content=201441621738&gclid=CI-RwZ_p7dQCFUa4wAodQ1MJUg


Default + General Drug + Painkillers + Stimulants + Sedatives
<iframe src="https://player.vimeo.com/video/240205224" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
 
Heroin
<iframe src="https://player.vimeo.com/video/240844000" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
 
Alcohol
<iframe src="https://player.vimeo.com/video/240844156" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
 
Cocaine
<iframe src="https://player.vimeo.com/video/240844853" width="640" height="360" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>