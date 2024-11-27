import Vue from 'vue'
import Vuex from 'vuex'
import getters from './getters'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        read_once: false,
        not_found: false,
        config: null,
        view: "loading",
        key: "",
        content: "",
        langtype: ""
    },
    mutations: {
        init(state) {
            state.not_found = state.read_once = false;
        },
        updateMode(state, payload) {
            state.read_once = payload.read_once;
        },
        updateNotFound(state, payload) {
            state.not_found = payload.not_found;
        },
        updateState(state, payload) {
            Object.assign(state, payload);
        }
    },
    getters
});
