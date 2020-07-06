import {Login} from "./widget/login.js";
import {Filters} from "./com/filter.js";


$(function() {
    new Filters().promise().then(function () {
        new Login({
            parent: "#Container",
        })
    })
});

