var $ = require("./jquery-3.1.0.js");

function loadSeason(season_key, success) {
    console.log("loading season: " + season_key);

    if (season_key == null) {
        season_key = "current";
    }
    path = "/api/seasons/" + season_key;

    $.get(path, function(data) {
        var season = JSON.parse(data);

        if (success) {
            success(season);
        }
    });

}

module.exports = loadSeason;
