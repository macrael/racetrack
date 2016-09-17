// play_types.js

var $ = require("./jquery-3.1.0.js");
var Handlebars = require("./handlebars-v4.0.5.js");
var api = require("./racetrack_api.js");

var play_types = {
    loadEditView: function(season, router) {

        console.log("loading playtypes");

        var edit_main_template = Handlebars.compile($("#edit-main").html());
        var edit_play_types_template = Handlebars.compile($("#edit-play-types").html());

        var new_content = edit_main_template({season_title: season["title"], 
                                                                                              bod: edit_play_types_template({play_types: season["play_types"]})})

        $("#content").html(new_content);
        $("#new-play-type-form").submit(function(e) {
            e.preventDefault();
            console.log("SUBMITIT");
            var newPlayType = {
                action: $("#new-play-type-action").val(),
                effect: Number($("#new-play-type-effect").val())
            };
            api.postObject(season["key"], "play_types", newPlayType, function (){
                console.log("succes post PT");
                router.reresolve();
            });
        });

        season.play_types.forEach(function(play_type, index) {
            console.log(play_type);
            var del_key = "#DEL" + play_type.key;
            del_key = del_key.replace(/:/, '\\:');
            $(del_key).click(function() {
                api.deleteObject(season["key"], "play_types", play_type.key, function() {
                    console.log("deletinPTg!");
                    router.reresolve();
                });
            });
        });

    }
}


module.exports = play_types;
