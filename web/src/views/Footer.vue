<template>
    <div class="row">
        <div class="col-md-12">
            <div class="footer">
                <p>
                    <a>Copyright&nbsp;&copy;&nbsp;2021&nbsp;-&nbsp;{{ year }}&nbsp;&nbsp;|&nbsp;&nbsp;</a>
                    <a href="https://beian.miit.gov.cn" title="备案号" target="_blank">{{ $store.state.config.beian.number }}</a>
                    <a v-for="footer in $store.state.config.footer" v-bind:key="footer.id">&nbsp;&nbsp;|&nbsp;&nbsp;<a :href="footer.url" target="_blank">{{ footer.text }}</a>
                    </a>
                </p>
                <p title="赞助">
                    <a v-b-modal.donate class="logo">
                        <img src="../assets/img/sponsor.svg" alt="打赏">
                    </a>
                </p>
            </div>
        </div>
        <b-modal id="donate" hide-footer lazy>
            <img src="../assets/img/sponsor.jpg" alt="赞赏">
        </b-modal>
    </div>
</template>

<script>
    export default {
        name: "Footer",
        data() {
            return {
                year: new Date().getFullYear(),
            }
        },
    }
</script>

<style scoped>
    .footer {
        font-size: .8em;
        text-align: center;
    }

    .footer p {
        margin: 1em;
    }

    .footer a:link, .footer a:visited {
        color: #000000;
    }

    #one-word {
        -webkit-touch-callout: none;
        -webkit-user-select: none;
        -khtml-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
        cursor: pointer;
    }

    #one-popover {
        font-family: Menlo, Monaco, "Andale Mono", "lucida console", "Courier New", monospace;
    }

    #donate img {
        width: 100%;
    }

    .logo img {
        height: 2em;
    }
    .logo svg {
        height: 2em;
    }
    .logo {
        margin: .8em;
    }
</style>
