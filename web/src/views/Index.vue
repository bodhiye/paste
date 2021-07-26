<template>
    <transition name="component-fade" mode="out-in">
        <component :is="view" v-bind="$data"></component>
    </transition>
</template>

<script>
    import stateMixins from "../mixins/stateMixin"
    import {mapGetters} from "vuex"
    import Form from "./Form"
    import Success from "./Success"
    import PasswordAuth from "./PasswordAuth"
    import PasteView from "./PasteView"
    import Loading from "./Loading"
    export default {
        name: "Index",
        mixins: [stateMixins],
        data() {
            return {}
        },
        computed: {
            ...mapGetters([
                "view",
                "langtype",
                "content"
            ])
        },
        watch: {
            "$route.params.key": function () {
                this.init();
            }
        },
        mounted() {
            this.init();
        },
        methods: {
            init() {
                this.$store.commit("init");
                if (this.$route.params.key === "") {
                    this.updateView("home");
                } else {
                    this.updateView("loading");
                    this.api.get(this.$store.getters.config.api.backend + this.$store.getters.config.api.path + '/' + this.$route.params.key, {
                        json: true
                    }).then(response => {
                        if (response.code === 200) {
                            this.updateView("paste_view");
                            this.updateContent(response.content);
                            this.updateLangtype(response.langtype === "plain" ? "plaintext" : response.langtype);
                        } else if (response.code === 401) {
                            this.updateView("password_auth");
                        } else if (response.code === 404 && this.$route.params.key.search("[a-zA-Z]{1}") !== -1) {
                            this.$store.commit("updateMode", {
                                read_once: true,
                            });
                            this.updateView("home");
                        } else {
                            this.$router.push("not_found");
                        }
                    });
                }
            },
        },
        components: {
            "home": Form,
            "success": Success,
            "password_auth": PasswordAuth,
            "paste_view": PasteView,
            "loading": Loading
        }
    }
</script>

<style scoped>
    .component-fade-enter-active, .component-fade-leave-active {
        transition: opacity .6s ease;
    }

    .component-fade-enter, .component-fade-leave-to {
        opacity: 0;
    }
</style>
