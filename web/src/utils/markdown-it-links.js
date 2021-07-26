
function markdownItLinkTarget(md, config) {
    config = config || {};

    const defaultRender = md.renderer.rules.link_open || this.defaultRender;
    const target = config.target || '_blank';

    md.renderer.rules.link_open = function (tokens, idx, options, env, self) {
        // If you are sure other plugins can't add `target` - drop check below
        const aIndex = tokens[idx].attrIndex('target');

        if (aIndex < 0) {
            tokens[idx].attrPush(['target', target]) // add new attribute
        } else {
            tokens[idx].attrs[aIndex][1] = target // replace value of existing attr
        }

        // pass token to default renderer.
        return defaultRender(tokens, idx, options, env, self)
    }
}

markdownItLinkTarget.defaultRender = function (tokens, idx, options, env, self) {
    return self.renderToken(tokens, idx, options)
};

module.exports = markdownItLinkTarget;
