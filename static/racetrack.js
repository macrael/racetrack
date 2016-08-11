// racetrack.js
// pulls down data from the webserver
// displayes the views

// function pathJoin(parts, sep){
//    var separator = sep || '/';
//    var replace   = new RegExp(separator+'{1,}', 'g');
//    return parts.join(separator).replace(replace, separator);
// }

var TheSeason;

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

// ---------- Routers ----------

var CurrentViewLoader;

function loadMainView(season) {

}

function loadEditView(season) {

}

function loadEditQueensView(season) {
    console.log("loadEditQueens");
    console.log(TheSeason)

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
        console.log(del_key);
        console.log($(del_key));
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
    router.on("/edit", function() {
        console.log("EDIT");
        CurrentViewLoader = loadEditView;
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