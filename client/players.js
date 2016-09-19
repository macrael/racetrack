// players.js

var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var calculator = require("./calculate.js");
var api = require("./racetrack_api.js");
var ui = require("./ui.js");

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
            bod: player_standings_template({players: sortedPlayers,
                current_ep_number: season["most_recent_episode_num"]
                }) });

        $("#content").html(playerStandings);
    },

    loadEditView: function (season, router, extra_data) {
        console.log("editingPlayers");
        
        var edit_main_template = Handlebars.compile($("#edit-main").html());
        var edit_players_template = Handlebars.compile($("#edit-players").html());
        var new_player_template = Handlebars.compile($("#new-player").html());
        var players_list_template = Handlebars.compile($("#players-list").html());

        if (!extra_data) {
            extra_data = {
                selected_queens: [],
                player_name: "Dummy",
                player_team: "Gay Idiots"
            }
        }

        var enable_create = extra_data.selected_queens.length == 3 && extra_data.selected_queens.find(function(queen) { return queen.winner });

        // add queens to players
        season.players.forEach(function(player) {
            var queens = []
            player.queen_keys.forEach(function(queen_key) {
                console.log("we check key", queen_key);
                queen = season.queens.find(function(q) {
                    return q.key == queen_key;
                });
                if (queen) {
                    var queen_copy = JSON.parse(JSON.stringify(queen));
                    if (queen_copy.key == player.winner_key) {
                        queen_copy.winner = true;
                    }
                    queens.push(queen_copy);
                } else {
                    console.log("DATA INTEGRIGTY BROKE!", queen_key);
                }
            });
            console.log("got queens", queens);
            player.queens = queens;
        });

        var edit_players = edit_main_template({
            bod: edit_players_template({
                new_player: new_player_template({
                    player_name: extra_data.player_name,
                    player_team: extra_data.player_team,
                    queens: season.queens,
                    selected_queens: extra_data.selected_queens,
                    create_enabled: enable_create
                }),
                old_players: players_list_template({ players: season.players})
                })
            });

        $("#content").html(edit_players);

        season.queens.forEach(function(queen, index) {
            var queen_button = "#QU" + queen.key;
            queen_button = queen_button.replace(/:/, '\\:');
            $(queen_button).click(function(e) {
                console.log("clicked Queen");
                e.preventDefault(); // why here and not elsewhere? any button in a form?
                queen.selected = !queen.selected;
                if (queen.selected) {
                    extra_data.selected_queens.push(queen);
                    if (extra_data.selected_queens.length > 3) {
                        extra_data.selected_queens[0].selected = false;
                        extra_data.selected_queens.shift();
                    }
                } else {
                    qi = extra_data.selected_queens.indexOf(queen);
                    extra_data.selected_queens.splice(qi, 1);
                }
                
                console.log(season);
                console.log(extra_data.selected_queens);
                extra_data.player_name = $("#player-name").val();
                extra_data.player_team = $("#team-name").val();

                players.loadEditView(season, router, extra_data);
            });
        });
        
        extra_data.selected_queens.forEach(function(selected_queen, index) {
            console.log("Selected Queen");
            var squeen_button = "#SQ" + selected_queen.key;
            squeen_button = squeen_button.replace(/:/, '\\:');
            $(squeen_button).click(function(e) {
                console.log("Picking Queen");
                e.preventDefault();
                
                extra_data.selected_queens.forEach(function(clear_queen, index) {
                    if (clear_queen === selected_queen) {
                        clear_queen.winner = !clear_queen.winner;
                        console.log("gota winner", clear_queen.winner);
                    } else {
                        clear_queen.winner = false;
                    }
                });
                extra_data.player_name = $("#player-name").val();
                extra_data.player_team = $("#team-name").val();
                players.loadEditView(season, router, extra_data);
            });
        });

        ui.configureDeleteButtons(season["key"], "players", season.players, router);

        $("#new-player-form").submit(function(e) {
            e.preventDefault();

            new_player = {
                name: $("#player-name").val(),
                team_name: $("#team-name").val(),
                queen_keys: extra_data.selected_queens.map(function(queen) { return queen.key }),
                winner_key: extra_data.selected_queens.find(function(queen) { return queen.winner}).key
            }

            api.postObject(season["key"], "players", new_player, function() {
                router.reresolve();
            });

        });
    }

}

module.exports = players;
