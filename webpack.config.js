module.exports = {
  entry: './js-src/app.js',

  output: {
    filename: './app/index.js',
  },

  devtool: 'source-map',

  module: {
    rules: [
      {
        test:    /\.elm$/,
        exclude: [/elm-stuff/, /node_modules/],
        loader:  'elm-webpack-loader?verbose=true&warn=true',
      },
    ],

    noParse: /\.elm$/,
  },

};
