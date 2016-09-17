var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var calculator = require("./calculate.js");

var players = {
    loadView: function (season) {
        console.log("display the players");
        var display_main_template = Handlebars.compile($("#display-main").html());
        var player_standings_template = Handlebars.compile($("#player-standings").html());

        calculator.calculate(season);

        var sortedPlayers = season["players"].sort(function(a, b) {
            return b.score - a.score;
        });

        var playerStandings = display_main_template({season_title: season["title"],
            bod: player_standings_template({players: sortedPlayers}) });

        $("#content").html(playerStandings);
    }

}

module.exports = players;
