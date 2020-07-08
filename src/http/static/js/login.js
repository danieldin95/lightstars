import {Login} from "./widget/login.js";
import {Template} from "./com/template.js";


$(function() {
    new Template().promise().then(function () {
        new Login({
            parent: "#login",
        })
    })
});

