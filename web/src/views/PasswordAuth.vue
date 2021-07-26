<template>
    <b-row>
        <b-col md="4"></b-col>
        <b-col md="4">
            <b-form @submit.prevent="onSubmit">
                <b-form-group :label="$t('lang.auth.form.label')">
                    <b-form-input
                    type="password"
                    v-model="form.password"
                    :placeholder="flag ? '' : this.$t('lang.auth.form.placeholder')">
                </b-form-input>
                </b-form-group>
                <b-button type="submit" variant="primary">{{ $t('lang.auth.form.button') }}</b-button>
            </b-form>
        </b-col>
    </b-row>
</template>

<script>
    import stateMixin from "../mixins/stateMixin";
    export default {
        name: "PasswordAuth",
        mixins: [stateMixin],
        data() {
            return {
                flag: true,
                form: {
                    password: null,
                }
            }
        },
        methods: {
            onSubmit() {
                const sendUrl = `${this.$store.getters.config.api.backend}${this.$store.getters.config.api.path}/${this.$route.params.key}?password=${this.form.password}`;
                this.api.get(sendUrl, {
                    json: 'true'
                }).then(response => {
                    if (response.code === 200) {
                        this.updateContent(response.content);
                        this.updateLangtype(response.langtype === "plain" ? "plaintext" : response.langtype);
                        this.updateView("paste_view");
                    } else {
                        this.flag = false;
                        this.form.password = null;
                    }
                });
            }
        }
    }
</script>

<style scoped>

</style>
