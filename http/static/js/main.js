import {Api} from "./api/api.js";
import {Location} from "./com/location.js";
import {Filters} from "./com/filter.js";
import {Navigation} from "./widget/navigation.js";
import {Routes} from "./routes.js";


$(function() {
    let hyper = $('hyper');
    let container = "#Container";
    let host = Location.query('node');

    if (host === undefined) {
        // if host is null, using default.
        host = hyper.attr('default');
        Location.query('node', host);
    }
    Api.Host(host);

    let nav = new Navigation({
        parent: "nav",
        home: ".",
        container: container,
        name: hyper.attr('name'),
    });
    let filter = new Filters();
    let routes = new Routes({
        hyper: hyper,
        container: container,
        onchange: function (e) {
           nav.refresh();
        },
    });
});

