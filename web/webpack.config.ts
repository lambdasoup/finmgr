import * as webpack from 'webpack'
import * as path from 'path'

module.exports = function(env: any): webpack.Configuration {
  return {
    entry: './src/index.ts',
    output: {
      path: path.resolve(__dirname, 'build'),
      filename: 'bundle.js'
    },
    module: {
      loaders: [
        {
          test: [/\.elm$/],
          exclude: [/elm-stuff/, /node_modules/],
          use: [
            { loader: 'elm-hot-loader' },
            {
              loader: 'elm-webpack-loader',
              options: env && env.production ? {} : { debug: true, warn: true }
            }
          ]
        },
        { test: /\.ts$/, loader: 'ts-loader' }
      ]
    },
    resolve: {
      extensions: ['.js', '.ts', '.elm']
      }
  }
}
