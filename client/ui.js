
var $ = require("./jquery-3.1.0.js");
var api = require("./racetrack_api.js");

var ui = {
    
    configureDeleteButtons: function(season_key, objects_name, objects, router) {
        
        objects.forEach(function(object, index) {
            var del_key = "#DEL" + object.key;
            del_key = del_key.replace(/:/, '\\:');
            $(del_key).click(function() {
                //TODO: make it a double-click situation.

                api.deleteObject(season_key, objects_name, object.key, function() {
                    console.log("deleted " + object.key + " from " + objects_name);
                    router.reresolve();
                });
            });
        });
    }

}

module.exports = ui;
