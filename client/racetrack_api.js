// racetrack_api.js
// functions to interact with the go server backend

var $ = require("./jquery-3.1.0.js");

var api = {

    postObject: function(season_key, objects_name, object, success) {
        path = "/api/seasons/" + season_key + "/" + objects_name;
        console.log("POST " + path);
        
        $.post(path, JSON.stringify(object), function() {
            console.log("Finihed POST");
            if (success) { success(); }
            console.log("Did POST");
        });

    },
    deleteObject: function(season_key, objects_name, object_key, success) {
        
        path = "/api/seasons/" + season_key + "/" + objects_name + "/" + object_key;
        console.log("DELETE " + path);

        $.ajax({
            url: path,
            type: 'DELETE',
            success: function(result) {
                // Do something with the result
                success();
            }
        });

    }

}

module.exports = api;
