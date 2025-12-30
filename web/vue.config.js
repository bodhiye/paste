const cdn = {
    development: {
        css: [
            "https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css",
            "https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.6.0/css/bootstrap.min.css",
            "https://cdnjs.cloudflare.com/ajax/libs/bootstrap-vue/2.21.2/bootstrap-vue.min.css"
        ],
        js: [
            "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.1.0/highlight.min.js"
        ]
    },
    production: {
        css: [
            "https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css",
            "https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.6.0/css/bootstrap.min.css",
            "https://cdnjs.cloudflare.com/ajax/libs/bootstrap-vue/2.21.2/bootstrap-vue.min.css"
        ],
        js: [
            "https://cdnjs.cloudflare.com/ajax/libs/vue/2.6.9/vue.runtime.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/vue-router/3.5.2/vue-router.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/vuex/3.6.2/vuex.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/axios/0.21.1/axios.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/mermaid/8.11.0/mermaid.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.1.0/highlight.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/bootstrap-vue/2.21.2/bootstrap-vue.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/markdown-it/12.1.0/markdown-it.min.js",
            "https://cdn.jsdelivr.net/npm/@chenfengyuan/vue-qrcode@1.0.2/dist/vue-qrcode.min.js",
            "https://cdnjs.cloudflare.com/ajax/libs/vue-i18n/8.25.0/vue-i18n.min.js"
        ]
    }
};

const CompressionWebpackPlugin = require('compression-webpack-plugin')

module.exports = {
    publicPath: '/',
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
        if (process.env.NODE_ENV === 'production') {
            config.externals = {
                "vue": "Vue",
                "vuex": "Vuex",
                "vue-router": "VueRouter",
                "vue-i18n": "VueI18n",
                "axios": "axios",
                "mermaid": "mermaid",
                "highlight.js": "hljs",
                "bootstrap-vue": "BootstrapVue",
                "markdown-it": "markdownit",
                "@chenfengyuan/vue-qrcode": "VueQrcode"
            }
        }
    },
    chainWebpack: config => {
        config.plugin('html').tap(args => {
            if (process.env.NODE_ENV === 'development') {
                args[0].cdn = cdn.development
            }
            if (process.env.NODE_ENV === 'production') {
                args[0].cdn = cdn.production
            }
            return args
        })
    }
};
