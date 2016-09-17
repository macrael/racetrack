// racetrack.js
// pulls down data from the webserver
// displayes the views

// This seems a little wrong...
var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var Router = require("./router.js");

var TheSeason;

// !--- API -----


function postQueen(queen, success) {
    console.log("postQueen", queen);

    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/queens"

    $.post(path, JSON.stringify(queen), function() {
        console.log("Finihed new Queen Post");
        Router.reresolve();
        console.log("Did queen");
    });

}

function postPlayer(player, success) {
    console.log("postPlayer", player);

    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/players"

    $.post(path, JSON.stringify(player), function() {
        console.log("Finihed new Player Post");
        Router.reresolve();
        console.log("Did player");
    });
}

function postPlayType(playType, success) {
    console.log("postPlayType", playType);

    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/play_types"

    $.post(path, JSON.stringify(playType), function() {
        console.log("Finihed new PlayType Post");
        Router.reresolve();
    });

}

function postEpisode(episode, success) {
    console.log("postEpisode", episode);

    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/episodes"

    $.post(path, JSON.stringify(episode), function() {
        console.log("Finihed new Episode Post");
        Router.reresolve();
    });

}


function postPlay(play) {
    console.log("postplay:", play);
    
    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/plays";

    $.post(path, JSON.stringify(play), function() {
        console.log("Finihed new Play");
        Router.reresolve();
    });

}
// ---------- Routers ----------

function loadQueensView(season) {
    console.log("my first display view");
    var display_main_template = Handlebars.compile($("#display-main").html());
    var queen_standings_template = Handlebars.compile($("#queen-standings").html());

    calculateQueenScores(season);
    console.log(season.queens);

    var sorted_queens = season["queens"].sort(function(a, b) {
        return b.score - a.score;
    });

    var queen_standings = display_main_template({season_title: season["title"],
        bod: queen_standings_template({queens: sorted_queens}) });

    $("#content").html(queen_standings);
}

function loadEditEpisodesView(season) {
    console.log("loading Episodes View");

    var edit_main_template = Handlebars.compile($("#edit-main").html());
    var edit_episodes_template = Handlebars.compile($("#edit-episodes").html());
    var new_episode_template = Handlebars.compile($("#new-episode").html());

    var sorted_eps = season["episodes"].sort(function(a, b) {
        return a.number - b.number;
    });
    var max_ep_num = 0;
    if (sorted_eps.length != 0) {
        max_ep_num = sorted_eps[sorted_eps.length - 1].number
    }
    var edit_episodes = edit_main_template({season_title: season["title"],
            bod: edit_episodes_template({episodes: season["episodes"], 
                                      new_episode: new_episode_template({next_number: max_ep_num + 1})})});

    $("#content").html(edit_episodes);
    $("#new-episode-form").submit(function(e) {
        e.preventDefault();
        console.log("SUBMITIddT");
        var new_episode = {
            number: Number($("#new-episode-number").val())
        };
        postEpisode(new_episode);
    });

    season.episodes.forEach(function(episode, index) {
        console.log(episode);
        var del_key = "#DEL" + episode.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation. -- could just alert
            season_key = season["key"];
            episode_path = "/api/seasons/" + season_key + "/episodes/" + episode.key; //ugh, urls are the best ids,

            $.ajax({
                url: episode_path,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETEDee");
                    Router.reresolve();
                }
            });
        });

        var edit_key = "#EDIT" + episode.key;
        edit_key = edit_key.replace(/:/, '\\:');
        $(edit_key).click(function() {
            edit_path = Router.generate("episode.edit", {episode_key: episode.key});
            console.log("editpa: ", edit_path);
            Router.navigate(edit_path);
        });
    });

}

function loadEditEpisodeView(season, episode_key) {
    console.log("Edit Episode:", episode_key)

    var edit_main_template = Handlebars.compile($("#edit-main").html());
    var edit_episode_template = Handlebars.compile($("#edit-episode").html());
    var new_play_template = Handlebars.compile($("#new-play").html());

    var episode = season.episodes.find(function(ep) { return ep.key === episode_key });
    console.log("TheIP: ", episode);

    var plays = season.plays.filter(function(play) { return play.episode_key === episode.key });
    plays = plays.sort(function(a,b) { return a.timestamp - b.timestamp});

    // attach the queen and play_type
    plays.forEach(function(play) {
        play.queen = season.queens.find(function(queen) { return queen.key === play.queen_key });
        play.play_type = season.play_types.find(function(play_type) { return play_type.key === play.play_type_key });
    });

    console.log("PLAYS: ", plays);

    var selected_queen = season.queens.find(function(queen) { return queen.selected });
    var selected_play_type = season.play_types.find(function(play_type) { return play_type.selected });

    var edit_episode = edit_main_template({season_title: season["title"],
        bod: edit_episode_template({number: episode.number,
            plays: plays,
            new_play: new_play_template({queens: season.queens, 
                                         play_types: season.play_types,
                                         create_enabled: (selected_queen && selected_play_type)})
        })});

    $("#content").html(edit_episode);
    window.scrollTo(0,10000000);

    season.queens.forEach(function(queen, index) {
        var queen_button = "#QU" + queen.key;
        queen_button = queen_button.replace(/:/, '\\:');
        $(queen_button).click(function() {
            console.log("clicked Queen");
            season.queens.forEach(function(clear_queen, index) {
                if (clear_queen === queen) {
                    clear_queen.selected = !clear_queen.selected;
                } else {
                    clear_queen.selected = false;
                }
            });
            loadEditEpisodeView(season, episode_key);
        });
    });

    season.play_types.forEach(function(play_type, index) {
        var play_type_button = "#PT" + play_type.key;
        play_type_button = play_type_button.replace(/:/, '\\:');
        $(play_type_button).click(function() {
            console.log("clicked PlatTye");
            season.play_types.forEach(function(clear_play_type, index) {
                if (clear_play_type === play_type) {
                    clear_play_type.selected = !clear_play_type.selected;
                } else {
                    clear_play_type.selected = false;
                }
            });
            loadEditEpisodeView(season, episode_key);
        });
    });

    plays.forEach(function(play, index) {
        var del_key = "#DEL" + play.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation. -- could just alert
            season_key = season["key"];
            play_path = "/api/seasons/" + season_key + "/plays/" + play.key; //ugh, urls are the best ids,

            $.ajax({
                url: play_path,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETEDpee");
                    Router.reresolve();
                }
            });
        });

    });
    
    // setup the submit button
    $("#add-play").click(function() {
        console.log("submitting");
        var new_play = {
            queen_key: selected_queen.key,
            play_type_key: selected_play_type.key,
            episode_key: episode.key
        };
        postPlay(new_play);
    });
    
}

function loadEditPlaysView(season) {
    console.log("loading plays");

    var edit_main_template = Handlebars.compile($("#edit-main").html());
    var edit_play_types_template = Handlebars.compile($("#edit-play-types").html());

    var new_content = edit_main_template({season_title: season["title"], 
                                                                                          bod: edit_play_types_template({play_types: season["play_types"]})})

    $("#content").html(new_content);
    $("#new-play-type-form").submit(function(e) {
        e.preventDefault();
        console.log("SUBMITIT");
        var newPlayType = {
            action: $("#new-play-type-action").val(),
            effect: Number($("#new-play-type-effect").val())
        };
        postPlayType(newPlayType);

        //TODO: pull down new truth?...... really just need to add this queen to the list....
    });

    season.play_types.forEach(function(play_type, index) {
        console.log(play_type);
        var del_key = "#DEL" + play_type.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation.
            season_key = season["key"];
            play_type_path = "/api/seasons/" + season_key + "/play_types/" + play_type.key; //ugh, urls are the best ids,

            $.ajax({
                url: play_type_path,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETEDss");
                    Router.reresolve()
                }
            });
        });
    });

}

function loadEditQueensView(season) {
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
        postQueen(newQueen);

        //TODO: pull down new truth?...... really just need to add this queen to the list....
    });
    // hook up all the delete buttons
    season.queens.forEach(function(queen, index) {
        console.log(queen);
        var del_key = "#DEL" + queen.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation.
            season_key = season["key"];
            queen_path = "/api/seasons/" + season_key + "/queens/" + queen.key; //ugh, urls are the best ids,

            $.ajax({
                url: queen_path,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETED");
                    Router.reresolve();
                }
            });
        });
    });

}

function loadEditPlayersView(season, extra_data) {
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

            loadEditPlayersView(season, extra_data);
        });
    });
    
    extra_data.selected_queens.forEach(function(selected_queen, index) {
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
            loadEditPlayersView(season, extra_data);
        });
    });

    season.players.forEach(function(player, index) {
        console.log(player);
        var del_key = "#DEL" + player.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation.
            season_key = season["key"];
            player_path = "/api/seasons/" + season_key + "/players/" + player.key; //ugh, urls are the best ids,

            $.ajax({
                url: player_path,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETplED");
                    Router.reresolve();
                }
            });
        });
    });

    $("#new-player-form").submit(function(e) {
        e.preventDefault();

        new_player = {
            name: $("#player-name").val(),
            team_name: $("#team-name").val(),
            queen_keys: extra_data.selected_queens.map(function(queen) { return queen.key }),
            winner_key: extra_data.selected_queens.find(function(queen) { return queen.winner}).key
        }

        postPlayer(new_player);

    });
}

// ---------- End Routers ----------

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
