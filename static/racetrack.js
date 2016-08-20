// racetrack.js
// pulls down data from the webserver
// displayes the views

// function pathJoin(parts, sep){
//    var separator = sep || '/';
//    var replace   = new RegExp(separator+'{1,}', 'g');
//    return parts.join(separator).replace(replace, separator);
// }

var TheSeason;

// !--- API -----

function loadSeason(season_key, success) {
    // !!dev vs. prod
    console.log("loading season: " + season_key);

    hostname = RacetrackConfig["server"];
    if (season_key == null) {
        season_key = "current";
    }
    path = "/api/seasons/" + season_key;

    season_url = hostname + path;
    console.log("Get: " + season_url)

    $.get(season_url, function(data) {
        console.log("ewback");
        console.log(data);

        var season = JSON.parse(data);

        TheSeason = season
        CurrentViewLoader(season);
        if (success) {
            success(season);
        }
    });

}

function postQueen(queen, success) {
    console.log("postQueen", queen);

    hostname = RacetrackConfig["server"];
    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/queens"

    new_queen_url = hostname + path;

    $.post(new_queen_url, JSON.stringify(queen), function() {
        console.log("Finihed new Queen Post");
        loadSeason(season_key)
    });

}

function postPlayType(playType, success) {
    console.log("postPlayType", playType);

    hostname = RacetrackConfig["server"];
    season_key = TheSeason["key"]
    path = "/api/seasons/" + season_key + "/play_types"

    new_play_type_url = hostname + path;

    $.post(new_play_type_url, JSON.stringify(playType), function() {
        console.log("Finihed new PlayType Post");
        loadSeason(season_key)
    });

}

// ---------- Routers ----------

var CurrentViewLoader;

function loadMainView(season) {

}

function loadEditEpisodesView(season) {

}

function loadEditPlaysView(season) {
    console.log("loading plays");

    var edit_main_template = Handlebars.compile($("#edit-main").html());
    var edit_play_types_template = Handlebars.compile($("#edit-play-types").html());

    var new_content = edit_main_template({season_title: TheSeason["title"], 
                                                                                          bod: edit_play_types_template({play_types: TheSeason["play_types"]})})

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

    TheSeason.play_types.forEach(function(play_type, index) {
        console.log(play_type);
        var del_key = "#DEL" + play_type.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation.
            hostname = RacetrackConfig["server"];
            season_key = TheSeason["key"];
            play_type_url = "/api/seasons/" + season_key + "/play_types/" + play_type.key; //ugh, urls are the best ids,

            $.ajax({
                url: play_type_url,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETEDss");
                    loadSeason(season_key)
                }
            });
        });
    });

}

function loadEditQueensView(season) {
    console.log("loadEditQueens");
    console.log(TheSeason);

    var edit_main_template = Handlebars.compile($("#edit-main").html());
    var edit_queens_template = Handlebars.compile($("#edit-queens").html());
    var new_queen_template = Handlebars.compile($("#new-queen").html());
    var queens_list_template = Handlebars.compile($("#queens-list").html());

    var new_content = edit_main_template({season_title: TheSeason["title"], bod: edit_queens_template({new_queen: new_queen_template(), old_queens: queens_list_template({queens: TheSeason.queens})})});

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
    TheSeason.queens.forEach(function(queen, index) {
        console.log(queen);
        var del_key = "#DEL" + queen.key;
        del_key = del_key.replace(/:/, '\\:');
        $(del_key).click(function() {
            //TODO: make it a double-click situation.
            hostname = RacetrackConfig["server"];
            season_key = TheSeason["key"];
            queen_url = "/api/seasons/" + season_key + "/queens/" + queen.key; //ugh, urls are the best ids,

            $.ajax({
                url: queen_url,
                type: 'DELETE',
                success: function(result) {
                    // Do something with the result
                    console.log("DELTEETED");
                    loadSeason(season_key)
                }
            });
        });
    });

}

function configureRouter() {
    router = new Navigo(null, false);
    
    router.on("/$", function() {
        console.log("slashhhhh");
    });
    router.on("/edit/queens", function() {
        console.log("EDITQ!");
        CurrentViewLoader = loadEditQueensView;
    });
    router.on("/edit/play_types", function() {
        console.log("EDITP!");
        CurrentViewLoader = loadEditPlaysView;
    });
    router.on("/edit", function() {
        console.log("EDIT");
        CurrentViewLoader = loadEditEpisodesView;
    });
    router.on("/queen", function() {
        console.log("queeeeen");
    });
    router.on(function() {
        console.log("404");
    });
    router.resolve();
}

// ---------- End Routers ----------


function debugPress() {
    console.log("debug press")

    var newQueen = {
        "name": "Bob the Drag Queen"
    }

}

function main() {
    console.log("hello there");
    $("#debug_button").click(debugPress);

    configureRouter();
    loadSeason(null, function(season) {
        console.log("Loaded!");

        var title = season["title"];
        // $("#season_name").html(title);
        
        console.log("foot");

    });

    console.log(RacetrackConfig)
}

$('document').ready(main);