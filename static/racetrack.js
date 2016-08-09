// racetrack.js
// pulls down data from the webserver
// displayes the views

// function pathJoin(parts, sep){
//    var separator = sep || '/';
//    var replace   = new RegExp(separator+'{1,}', 'g');
//    return parts.join(separator).replace(replace, separator);
// }

// var TheSeason;

function loadSeason(season_key, success) {
	// !!dev vs. prod
	console.log("loading season");

	hostname = RacetrackConfig["server"];
	path = "/api/seasons/";
	if (season_key == null) {
		season_key = "current";
	}

	season_url = hostname + path + season_key;

	$.get(season_url, function(data) {
		console.log("ewback");
		console.log(data);

		var season = JSON.parse(data);

		success(season)
	});

}