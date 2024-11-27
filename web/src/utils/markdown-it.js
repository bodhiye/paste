/* Author: Ryan Lee(ryanlee2014)
 * Created at 08/18/2019
 * MarkdownIt Instance
 */

const uslug = require("uslug")
const uslugify = s => uslug(s)

function Instance(key = "", problem_id = "") {
    const md = require("markdown-it")({
        html: true,
        linkify: true,
        typographer: true,
        breaks: true
    })
    const mh = require("markdown-it-highlightjs")
    const mk = require("@ryanlee2014/markdown-it-katex")
    const ma = require("markdown-it-anchor").default
    md.use(mk)
    md.use(mh)
    md.use(ma, {
        slugify: uslugify
    })
    md.use(require('markdown-it-task-checkbox'), {
        disabled: true,
        divWrap: false,
        divClass: 'checkbox',
        idPrefix: 'cbx_',
        ulClass: 'task-list',
        liClass: 'task-list-item'
    })
    md.use(require("./markdown-it-links"))
    md.use(require("./markdown-it-mermaid").default.default)

    const markdownPack = (html) => {
        // return `<div class="markdown-body">${html}</div>`
        return html;
    }

    const preToSegment = (html) => {
        return html.replace(/<pre>[\s\S]+?<\/pre>/g, `<div class='ui segment'>
    <div class="ui top attached label"><a class="copy context">Copy</a></div>$&</div>`)
    }

    const _render = md.render

    md.render = function () {
        return markdownPack(_render.apply(md, arguments))
    }

    md.renderRaw = function () {
        return preToSegment(md.renderInline(...arguments))
    }

    return Object.assign(md, {key, problem_id})
}

const md = Instance()
md.newInstance = Instance
export default md
