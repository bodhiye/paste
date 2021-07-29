<template>
    <b-row>
        <b-col md="1" lg="2"></b-col>
        <b-col md="10" lg="8">
            <b-form @submit.prevent="onSubmit">
                <b-row>
                    <b-col md="7" lg="5">
                        <b-form-group>
                            <b-input-group :prepend="$t('lang.form.input[0].prepend')">
                                <b-form-select v-model="form.langtype">
                                    <option value="plain">{{ $t('lang.form.select.plain') }}</option>
                                    <option value="bash">Bash</option>
                                    <option value="cpp">C/C++</option>
                                    <option value="go">Go</option>
                                    <option value="java">Java</option>
                                    <option value="json">JSON</option>
                                    <option value="python">Python</option>
                                    <option value="markdown">Markdown</option>
                                </b-form-select>
                            </b-input-group>
                        </b-form-group>
                        <b-form-group>
                            <b-input-group :prepend="$t('lang.form.input[1].prepend')">
                                <b-form-input type="password" autocomplete="off" v-model="form.password" :placeholder="$t('lang.form.input[1].placeholder')">
                                </b-form-input>
                            </b-input-group>
                        </b-form-group>
                        <b-form-group>
                            <b-input-group :prepend="$t('lang.form.input[2].prepend')">
                                <b-form-select v-model="form.expireDate">
                                    <option value="none">{{ $t('lang.form.select.none') }}</option>
                                    <option value="hour">{{ $t('lang.form.select.hour') }}</option>
                                    <option value="day">{{ $t('lang.form.select.day') }}</option>
                                    <option value="week">{{ $t('lang.form.select.week') }}</option>
                                    <option value="month">{{ $t('lang.form.select.month') }}</option>
                                    <option value="year">{{ $t('lang.form.select.year') }}</option>
                                </b-form-select>
                            </b-input-group>
                        </b-form-group>
                    </b-col>
                </b-row>
                <b-row>
                    <b-col md="12">
                        <b-form-group>
                            <b-form-textarea v-model="form.content" rows="12" 
                            :placeholder="$t('lang.form.textarea.placeholder.' + 
                            ($store.state.read_once ? 'read_once' : 'code'))" 
                            required no-resize maxlength="100000" oninvalid="this.setCustomValidity('show me the code')">
                            </b-form-textarea>
                        </b-form-group>
                        <b-form-group>
                            <b-checkbox-group switches>
                                <b-button type="submit" :variant="$store.state.read_once ? 'dark' : 'primary'" style="margin-right: .65em">
                                    {{ $t('lang.form.submit') }}
                                </b-button>
                                <b-form-checkbox v-model="read_once" v-show="!$store.state.read_once" switch>
                                    {{ $t('lang.form.checkbox') }}
                                </b-form-checkbox>
                            </b-checkbox-group>
                        </b-form-group>
                    </b-col>
                </b-row>
            </b-form>
        </b-col>
        <b-col md="1" lg="2"></b-col>
    </b-row>
</template>

<script>
    import stateMixins from "../mixins/stateMixin";
    export default {
        name: "Form",
        mixins: [stateMixins],
        data() {
            return {
                form: {
                    langtype: 'plain',
                    content: null,
                    password: null,
                    expireDate: 'none'
                },
                read_once: []
            }
        },
        methods: {
            onSubmit() {
                let key = "";
                if (this.$route.params.key !== '') {
                    key = this.$route.params.key;
                } else if (this.read_once.length > 0) {
                    key = "/once"
                }

                if (this.form.expireDate === "none") {
                    this.form.expireDate = 0
                } else if (this.form.expireDate === 'hour') {
                    this.form.expireDate = 60*60
                } else if (this.form.expireDate === "day") {
                    this.form.expireDate = 24*60*60
                } else if (this.form.expireDate === "week") {
                    this.form.expireDate = 7*24*60*60
                } else if (this.form.expireDate === "month") {
                    this.form.form.expireDate = 30*24*60*60
                } else if (this.form.expireDate === "year") {
                    this.form.expireDate = 365*24*60*60
                }

                const sendArgs = [`${this.$store.getters.config.api.backend}${this.$store.getters.config.api.path}${key}`, this.form];
                const sendFunc = this.api.post;
                sendFunc(...sendArgs).then(response => {
                    if (response.code === 201) {
                        this.updateView("success");
                        this.updateKey(response.key);
                    }
                });
            }
        }
    }
</script>

<style scoped>

</style>
