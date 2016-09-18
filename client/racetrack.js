// racetrack.js
// pulls down data from the webserver
// displayes the views

// This seems a little wrong...
var $ = require("./jquery-3.1.0.js");
var Router = require("./router.js");

var TheSeason;

function debugPress() {
    console.log("debug press")

    var newQueen = {
        "name": "Bob the Drag Queen"
    }
    $.get("/debug", function(data) {
        console.log(data);
    });

}

function main() {
    console.log("hello there");
    $("#debug_button").click(debugPress);

    Router.resolve();
}

$('document').ready(main);
