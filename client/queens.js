// queens.js

var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var calculator = require("./calculate.js");
var api = require("./racetrack_api.js");
var ui = require("./ui.js");

var queens = {
    loadView: function (season) {
        console.log("my first display view");
        var display_main_template = Handlebars.compile($("#display-main").html());
        var queen_standings_template = Handlebars.compile($("#queen-standings").html());

        calculator.calculate(season);

        var sorted_queens = season["queens"].sort(function(a, b) {
            return b.score - a.score;
        });

        var queen_standings = display_main_template({season_title: season["title"],
            bod: queen_standings_template({queens: sorted_queens}) });

        $("#content").html(queen_standings);
    },

    loadEditView: function (season, router) {
        console.log("loadEditQueens");
        console.log(season);

        var edit_main_template = Handlebars.compile($("#edit-main").html());
        var edit_queens_template = Handlebars.compile($("#edit-queens").html());
        var new_queen_template = Handlebars.compile($("#new-queen").html());
        var queens_list_template = Handlebars.compile($("#queens-list").html());

        var new_content = edit_main_template({season_title: season["title"], bod: edit_queens_template({new_queen: new_queen_template(), old_queens: queens_list_template({queens: season.queens})})});

        $("#content").html(new_content);
        $("#new-queen-form").submit(function(e) {
            e.preventDefault();
            console.log("SUBMIT");
            var newQueen = {
                name: $("input#name").val()
            };
            api.postObject(season["key"], "queens", newQueen, function () {
                console.log("success post QU");
                router.reresolve();
            });
        });

        ui.configureDeleteButtons(season["key"], "queens", season.queens, router);

    }
}

module.exports = queens;
