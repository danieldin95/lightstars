import {Api} from "./api/api.js";
import {Location} from "./com/location.js";
import {Template} from "./com/template.js";
import {Navigation} from "./widget/navigation.js";
import {Routes} from "./routes.js";


$(function() {
    let hyper = $('hyper');
    let host = Location.query('node');

    if (!host) {
        // if host is null, using default.
        host = hyper.attr('default');
        Location.query('node', host);
    }
    Api.Host(host);

    new Template().promise().then(function () {
        let nav = new Navigation({
            parent: "#Navigation",
            home: ".",
            container: "#Container",
            name: hyper.attr('name'),
        });
        let rte = new Routes({
            hyper: hyper,
            container: "#Container",
            onchange: function (e) {
                nav.refresh();
            },
        });
    });
});

