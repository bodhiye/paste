import mermaid from "mermaid";

const o = function (e) {
    function r(n) {
        if (t[n]) return t[n].exports;
        var a = t[n] = {i: n, l: !1, exports: {}};
        return e[n].call(a.exports, a, a.exports, r), a.l = !0, a.exports
    }

    var t = {};
    return r.m = e, r.c = t, r.d = function (e, t, n) {
        r.o(e, t) || Object.defineProperty(e, t, {configurable: !1, enumerable: !0, get: n})
    }, r.n = function (e) {
        var t = e && e.__esModule ? function () {
            return e.default
        } : function () {
            return e
        };
        return r.d(t, "a", t), t
    }, r.o = function (e, r) {
        return Object.prototype.hasOwnProperty.call(e, r)
    }, r.p = "", r(r.s = 0)
}([function (e, r, t) {
    "use strict";
    Object.defineProperty(r, "__esModule", {value: !0});
    var n = t(1), a = function (e) {
        return e && e.__esModule ? e : {default: e}
    }(n), u = function (e) {
        try {
            return a.default.parse(e), '<div class="mermaid">' + e + "</div>"
        } catch (e) {
            var r = e.str;
            e.hash;
            return "<pre>" + r + "</pre>"
        }
    }, i = function (e) {
        e.mermaid = a.default, a.default.loadPreferences = function (e) {
            var r = e.get("mermaid-theme");
            void 0 === r && (r = "default");
            var t = e.get("gantt-axis-format");
            return void 0 === t && (t = "%Y-%m-%d"), a.default.initialize({
                theme: r,
                gantt: {
                    axisFormatter: [[t, function (e) {
                        return 1 === e.getDay()
                    }]]
                },
                startOnLoad: false
            }), {"mermaid-theme": r, "gantt-axis-format": t}
        };
        var r = e.renderer.rules.fence.bind(e.renderer.rules);
        e.renderer.rules.fence = function (e, t, n, a, i) {
            var o = e[t], f = o.content.trim();
            if ("mermaid" === o.info) return u(f);
            var c = f.split(/\n/)[0].trim();
            return "gantt" === c || "sequenceDiagram" === c || c.match(/^graph (?:TB|BT|RL|LR|TD);?$/) ? u(f) : r(e, t, n, a, i)
        }
    };
    r.default = i
// eslint-disable-next-line no-unused-vars
}, function (e, r) {
    e.exports = mermaid
}]);
export default o;
//# sourceMappingURL=index.js.map
