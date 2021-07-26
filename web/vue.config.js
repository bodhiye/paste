const cdn = {
    css: [
        "https://cdn.staticfile.org/github-markdown-css/4.0.0/github-markdown.min.css",
        "https://cdn.staticfile.org/twitter-bootstrap/4.6.0/css/bootstrap.min.css",
    ],
    js: [
        "https://cdn.staticfile.org/highlight.js/11.1.0/highlight.min.js"
    ]
};

const CompressionWebpackPlugin = require('compression-webpack-plugin')

module.exports = {
    publicPath: './',
    outputDir: 'paste',
    assetsDir: 'static',
    lintOnSave: process.env.NODE_ENV !== 'production',
    productionSourceMap: false,
    devServer: {
        port: process.env.VUE_APP_CLI_PORT,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            [process.env.VUE_APP_BASE_API]: { //需要代理的路径
                target: `${process.env.VUE_APP_BASE_PATH}:${process.env.VUE_APP_SERVER_PORT}`,
                changeOrigin: true,
                pathRewrite: { // 修改路径数据
                    ['^' + process.env.VUE_APP_BASE_API]: '/'
                }
            }
        }
    },
    chainWebpack: config => {
        config.plugin('html').tap(args => {
            args[0].cdn = cdn
            return args
        })
    },
    configureWebpack: config => {
        const productionGzipExtensions = ['html', 'js', 'css'];
        config.plugins.push(
            new CompressionWebpackPlugin({
                algorithm: 'gzip',
                test: new RegExp('\\.(' + productionGzipExtensions.join('|') + ')$'),
                threshold: 10240,
                minRatio: 0.8,
                deleteOriginalAssets: false
            })
        )
    }
};
