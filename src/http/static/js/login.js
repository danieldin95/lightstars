import {Login} from "./widget/login.js";
import {I18N} from "./com/i18n.js";
import {Template} from "./com/template.js";


$(function() {
    I18N.promise().then(function () {
        let tmpl = new Template();
        let login = new Login({
            parent: "#login",
        })
    })
});

