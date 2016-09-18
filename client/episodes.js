// episodes.js

var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var api = require("./racetrack_api.js");

var episodes = {

    loadEditView: function(season, router) {
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
            api.postObject(season["key"], "episodes", new_episode, function() {
                console.log("succ post ep");
                router.reresolve();
            });
        });

        season.episodes.forEach(function(episode, index) {
            console.log(episode);
            var del_key = "#DEL" + episode.key;
            del_key = del_key.replace(/:/, '\\:');
            $(del_key).click(function() {
                api.deleteObject(season["key"], "episodes", episode.key, function() {
                    console.log("DELTEETEDee");
                    router.reresolve();
                });
            });

            var edit_key = "#EDIT" + episode.key;
            edit_key = edit_key.replace(/:/, '\\:');
            $(edit_key).click(function() {
                edit_path = router.generate("episode.edit", {episode_key: episode.key});
                console.log("editpa: ", edit_path);
                router.navigate(edit_path);
            });
        });

    },

    loadEditEpisodeView: function(season, episode_key, router) {
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
                episodes.loadEditEpisodeView(season, episode_key, router);
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
                episodes.loadEditEpisodeView(season, episode_key, router);
            });
        });

        plays.forEach(function(play, index) {
            var del_key = "#DEL" + play.key;
            del_key = del_key.replace(/:/, '\\:');
            $(del_key).click(function() {
                api.deleteObject(season["key"], "plays", play.key, function() {
                    console.log("DELTEETEDpee");
                    router.reresolve();
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
            api.postObject(season["key"], "plays", new_play, function() {
                console.log("posted that play!");
                router.reresolve();
            });

        });
        
    }

}

module.exports = episodes;
