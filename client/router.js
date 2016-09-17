// router.js
// create the router and configure the routes

var Navigo = require("./navigo.js");

var loadSeason = require("./load_season.js");
var players = require("./players.js");
var queens = require("./queens.js");

var router = new Navigo(null, false);

router.reresolve = function() {
    // This is a hack, Navigo's regular behavior when resolve is called and
    // you are already at that location is to do nothing. This works around
    // that so that it will actually run the callback again in this case. 
    this._lastRouteResolved = null;
    this.resolve();
}

router.on("/$", function() {
    console.log("slashhhhh");
    loadSeason(null, function(season) {
        players.loadView(season);
    });
});
router.on("/edit/queens", function() {
    console.log("EDITQ!");
    loadSeason(null, function(season) {
        queens.loadEditView(season, router); // I'm not a fan of passing the router down, but not sure how to let the views trigger a refresh otherwise. 
    });
});
router.on("/queens", function() {
    console.log("SHOW! Q");
    loadSeason(null, function(season) {
        queens.loadView(season);
    });
});
router.on("/edit/play_types", function() {
    console.log("EDITP!");
    loadSeason(null, function(season) {
        loadEditPlaysView(season);
    });
});
router.on({"/edit/episodes/:episode_key": { 
    as:"episode.edit", 
    uses: function(params) {
        console.log("edit taht episode: ", params.episode_key);
        loadSeason(null, function(season) {
            loadEditEpisodeView(season, params.episode_key);
        });
    }
}});
router.on("/edit/episodes", function() {
    console.log("EDIT");
    loadSeason(null, function(season) {
        loadEditEpisodesView(season);
    });
});
router.on("/edit/players", function() {
    console.log("EDITP");
    loadSeason(null, function(season) {
        players.loadEditView(season, router);
    });
});
router.on("/edit$", function() {
    console.log("EDIT");
    loadSeason(null, function(season) {
        loadEditEpisodesView(season);
    });
});
router.on(function() {
    console.log("404");
});

module.exports = router;
