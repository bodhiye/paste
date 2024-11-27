import Vue from "vue"
import Router from "vue-router"
import Index from "../views/Index";

const emptyFunc = () => {};
const warn = process.env.NODE_ENV !== "production" ? (console && console.warn || emptyFunc) : emptyFunc;
const originalPush = Router.prototype.push;
Router.prototype.push = function push(location, onResolve, onReject) {
    if (onResolve || onReject) return originalPush.call(this, location, onResolve, onReject);
    return originalPush.call(this, location).catch(warn);
};

Vue.use(Router);

const router = new Router({
    mode: "history",
    base: "/",
    routes: [
        {
            path: "/:key(0{0}|[0-9a-zA-Z]{10})",
            name: "index",
            component: Index
        },
        {
            path: "/not_found",
            name: "NotFound",
            component: () => import("../views/NotFound")
        },
        {
            path: "*",
            redirect: "/not_found"
        }
    ]
})

router.beforeEach(async (to, from, next) => {
    if (to.path) {
        if (window._hmt) {
            window._hmt.push(['_trackPageview', to.fullPath])
        }
    }
    next()
})
export default router
