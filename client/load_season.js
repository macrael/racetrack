var $ = require("./jquery-3.1.0.js");

function loadSeason(season_key, success) {
    // !!dev vs. prod
    console.log("loading season: " + season_key);

    if (season_key == null) {
        season_key = "current";
    }
    path = "/api/seasons/" + season_key;

    console.log("Get: " + path)

    $.get(path, function(data) {
        console.log("ewback");
        console.log(data);

        var season = JSON.parse(data);

        TheSeason = season
        if (success) {
            success(season);
        }
    });

}

module.exports = loadSeason;
