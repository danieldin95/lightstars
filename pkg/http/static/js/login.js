import {LoginWid} from "./widget/login.js";
import {I18N} from "./lib/i18n.js";
import {Template} from "./lib/template.js";


$(function() {
    $.removeCookie('token', { path: '/' });
    I18N.promise().then(function () {
        new Template();
        new LoginWid({
            parent: "#login",
        })
    })
});
