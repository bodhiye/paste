export default {
    methods: {
        _baseUpdate(object) {
            this.$store.commit("updateState", object);
        },
        updateContent(content) {
            this._baseUpdate({ content });
        },
        updateLangtype(langtype) {
            this._baseUpdate({ langtype });
        },
        updateKey(key) {
            this._baseUpdate({ key });
        },
        updateView(view) {
            this._baseUpdate({ view });
        }
    }
}
