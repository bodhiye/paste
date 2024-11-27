<template>
    <b-row>
        <b-col md="1"></b-col>
        <b-col md="10" id="view">
            <div>
                <b-card no-body>
                    <b-card-header>
                        <b-row>
                            <b-col md="6">
                                <div>
                                    <a>{{ linesCount }} {{ $t('lang.view.lines') }}</a>
                                    <a>&nbsp;|&nbsp;</a>
                                    <a>{{ $t('lang.view.lang.' + langtype) }}</a>
                                </div>
                            </b-col>
                            <b-col md="6" style="text-align: right;">
                                <b-check-group switches>
                                    <b-checkbox v-model="raw" v-if="langtype === 'markdown'">{{ $t('lang.view.raw') }}</b-checkbox>
                                    <b-link id="clipboard-btn" :data-clipboard-text="content">
                                        {{ $t('lang.view.copy') }}
                                    </b-link>
                                    <b-tooltip show target="clipboard-btn" placement="bottomleft">
                                        {{ $t('lang.view.tooltip.' + (copy_btn_status > 0 ? 'success' :
                                        (copy_btn_status === 0 ? 'click' : 'fail'))) }}
                                    </b-tooltip>
                                </b-check-group>
                            </b-col>
                        </b-row>
                    </b-card-header>
                    <div ref="hljs">
                        <b-card-body style="padding-bottom: 0" v-if="langtype !== 'markdown' || raw.length === 1">
                            <pre><code v-bind:class="'line-numbers ' + langtype" v-text="this.content"></code></pre>
                        </b-card-body>
                        <b-card-body style="padding-bottom: 0" v-else>
                            <div class="markdown-body">
                                <div v-html="markdown.render(content)"></div>
                            </div>
                        </b-card-body>
                    </div>
                </b-card>
            </div>
        </b-col>
        <b-col md="1"></b-col>
    </b-row>
</template>

<script>
    import lineNumbersBlock from '../utils/highlightjs-line-numbers'
    import '../style/github-gist.css'
    import '../style/highlightjs-line-numbers.css'
    import { mapGetters } from "vuex"
    export default {
        name: "PasteView",
        data() {
            return {
                copy_btn_status: 0,
                raw: []
            }
        },
        mounted() {
            let clipboard = new this.clipboard('#clipboard-btn');
            let cur = this;
            clipboard.on('success', function () {
                cur.copy_btn_status = 1;
                window.getSelection().removeAllRanges();
                window.setTimeout(function () {
                    cur.copy_btn_status = 0;
                }, 2000);
            });
            clipboard.on('error', function () {
                cur.copy_btn_status = -1;
                window.setTimeout(function () {
                    cur.copy_btn_status = 0;
                }, 2000);
            });
            this.initMermaid();
            this.renderHljs(this.$refs.hljs);
            this.markdownBindHook();
        },
        computed: {
            linesCount: function () {
                let BREAK_LINE_REGEXP = /\r\n|\r|\n/g;
                return (this.content.trim().match(BREAK_LINE_REGEXP) || []).length + 1;
            },
            ...mapGetters([
                "langtype",
                "content"
            ])
        },
        methods: {
            renderHljs: function (el) {
                this.$nextTick(() => {
                    let blocks = el.querySelectorAll('pre code');
                    if (document.querySelectorAll('.hljs').length === 0) {
                        blocks.forEach(function (block) {
                            window.hljs.highlightBlock(block);
                            if (block.getAttribute('class').split(' ').indexOf('line-numbers') > -1) {
                                lineNumbersBlock(block, {
                                    singleLine: true
                                });
                            }
                        });
                    }
                })
            },
            initMermaid: function () {
                this.$nextTick(() => {
                    if (this.langtype === "markdown") {
                        import("mermaid").then(mermaid => {
                            document.querySelectorAll(".mermaid").forEach(v => {
                                mermaid.init(undefined, v);
                            });
                        })
                    }
                })
            },
            markdownBindHook: function () {
                const _render = this.markdown.render;
                const that = this;
                this.markdown.render = function () {
                    that.initMermaid();
                    return _render.apply(this, arguments);
                }
            }
        },
        watch: {
            raw: function () {
                this.renderHljs(document);
            },
            content: function () {
                this.initMermaid();
            }
        }
    }
</script>
<style>
    .markdown-body code {
        color: black!important
    }
</style>
<style scoped>
    #view {
        margin-top: 2.5em;
    }

    .markdown-body {
        box-sizing: border-box;
        min-width: 200px;
        max-width: 980px;
        margin: 0 auto;
        padding: 45px;
    }

    .markdown-body pre {
        padding-left: 1em;
    }

    @media (max-width: 767px) {
        .markdown-body {
            padding: 15px;
        }
    }
</style>
