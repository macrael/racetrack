module.exports = {
    context: __dirname + "/client",
    entry: "./racetrack.js",
    output: {
        path: __dirname + "/static",
        filename: "racetrack.js"
    }
    // module: {
    //     loaders: [
    //         { test: /\.css$/, loader: "style!css" }
    //     ]
    // }
};
