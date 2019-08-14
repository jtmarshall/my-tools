const webpack = require('webpack');
const config = {
    entry: __dirname + '/js/index.js',
    module: {
        rules: [
            {
                test: /\.jsx?/,
                exclude: /node_modules/,
                loader: 'babel-loader',
                options: {
                    presets: ['@babel/preset-env',
                        '@babel/react', {
                            'plugins': ['@babel/plugin-proposal-class-properties']
                        }]
                }
            }
        ]
    },
    output: {
        path: __dirname + '/dist',
        filename: 'bundle.js',
    },
    resolve: {
        extensions: ['.js', '.jsx', '.css']
    },
};
module.exports = config;