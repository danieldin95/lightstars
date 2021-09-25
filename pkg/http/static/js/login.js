import {Login} from "./widget/login.js";
import {I18N} from "./lib/i18n.js";
import {Template} from "./lib/template.js";


$(function() {
    $.removeCookie('token', { path: '/' });
    I18N.promise().then(function () {
        let tmpl = new Template();
        let login = new Login({
            parent: "#login",
        })
    })
});

