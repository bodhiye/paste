<template>
    <b-row>
        <b-col md="2"></b-col>
        <b-col md="8" id="success_fixed">
            <div class="jumbotron">
                <h2>
                    {{ $t('lang.success.h2') }}
                </h2>
                <ul>
                    <li>{{ $t('lang.success.ul.li[1].browser') }}
                        <a v-b-tooltip.hover="$t('lang.success.ul.li[1].tooltip')"
                           :href="base_url + key"
                           target="_blank">
                            {{ base_url + key }}
                        </a>&nbsp;<b-badge
                                variant="info"
                                class="badge-fixed"
                                :data-clipboard-text="base_url + key"
                                href="#">
                            {{ $t('lang.success.badge.' +
                            (copy_btn_status > 0 ? 'success' : (copy_btn_status === 0 ?  'copy' : 'fail')))  }}
                        </b-badge>
                    </li>
                    <li>
                        <a>{{ $t('lang.success.ul.li[2].scan_qr_code') }}</a>
                        <div class="text-left">
                            <QRCode :value="this.base_url + this.key" :options="{ color: { dark: '#0074d9' }, width: 111 }"></QRCode>
                        </div>
                    </li>
                </ul>
                <p>
                    <b-button @click.prevent="goHome" variant="primary">{{ $t('lang.success.p[0].button') }}</b-button>
                </p>
            </div>
        </b-col>
        <b-col md="2"></b-col>
        <b-popover
                :show.sync="popover_show"
                target="nav_input"
                placement="bottomright">
            <a v-html="$t('lang.success.popover.text')"></a>
        </b-popover>
    </b-row>
</template>

<script>
    import { mapGetters } from "vuex";
    import stateMixin from "../mixins/stateMixin";
    export default {
        name: "Success",
        mixins: [stateMixin],
        data() {
            return {
                base_url: location.origin + '/',
                copy_btn_status: 0,
                popover_show: false,
            }
        },
        computed: {
            ...mapGetters([
                "key"
            ])
        },
        mounted() {
            let clipboard = new this.clipboard('.badge-fixed');
            let cur = this;
            clipboard.on('success', function() {
                cur.copy_btn_status = 1;
                window.setTimeout(function () {
                    cur.copy_btn_status = 0;
                }, 2000);
            });
            clipboard.on('error', function() {
                cur.copy_btn_status = -1;
                window.setTimeout(function () {
                    cur.copy_btn_status = 0;
                }, 2000);
            });
        },
        methods: {
            goHome() {
                if (this.$route.params.key !== '') {
                    this.$router.push('/');
                } else {
                    this.updateView("home");
                }
            }
        }
    }
</script>

<style scoped>
    #success_fixed {
        margin-top: 1.375em;
        margin-bottom: 1.375em;
    }

    .badge-fixed {
        position: relative;
        bottom: .2em;
    }
</style>
