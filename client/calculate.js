// calculate.js
//
// just some functions that take the season object in and annotate it with calculated numbers. 


// ugh, maybe I want to do this at the start? Maybe the API should return it??
// Question is having the dict more useful at other times. As is might be OK
// Probably need to refactor this around
var _seasonDict;
function seasonDict(season) {
    if (_seasonDict) {
        return _seasonDict;
    }
    _seasonDict = {};
    for (var key in season) {
        if (season[key].constructor === Array) {
            _seasonDict[key] = {};
            season[key].forEach(function(item, index) {
               var itemKey = item["key"];
               if (itemKey == null) {
                   console.log("Woah, no key here... not sure what to do.");
                   throw "tried to be too clever, apparently"
               }
               _seasonDict[key][itemKey] = item;
            });
        } else {
            _seasonDict[key] = season[key];
        }
    }
    return _seasonDict;
}


// inserts "score" into all the queen objects
function calculateQueenScores(season) {

    var seasonD = seasonDict(season);

    for (var queenKey in seasonD["queens"]) {
        var queen = seasonD["queens"][queenKey];
        queen.score = 0;
    }

    for (var playKey in seasonD["plays"]) {
        var play = seasonD["plays"][playKey];
        var queen = seasonD["queens"][play["queen_key"]];
        var playType = seasonD["play_types"][play["play_type_key"]];

        queen.score = queen.score + playType.effect
    }

    season["queens"].forEach( function(queen, index) {
        dq = seasonD["queens"][queen["key"]];
        queen.score = dq.score;
    });

}

// this requires the queen scores to have been calculated already...
function calculatePlayerScores(season) {
    var seasonD = seasonDict(season);

    for (var playerKey in seasonD["players"]) {
        var player = seasonD["players"][playerKey];
        player.score = player["queen_keys"].map(function(queen_key) {
            return seasonD["queens"][queen_key];
        }).reduce(function(collector, queen) {
            return collector + queen.score;
        }, 0);
    }

    season["players"].forEach(function(player, index) {
        dp = seasonD["players"][player["key"]];
        player.score = dp.score;
    });
    
}

function calculateScores(season) {
    calculateQueenScores(season);
    calculatePlayerScores(season);
}


var calculator = {
    calculate: calculateScores
};

module.exports = calculator;


