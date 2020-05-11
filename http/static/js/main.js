import {Filters} from "./ctl/filter.js";
import {Navigation} from "./widget/navigation.js";
import {Routes} from "./routes.js";
import {Location} from "./com/location.js";
import {Api} from "./api/api.js";

console.log('init')

$(function() {
    let hyper = $('hyper');
    let container = "#container";
    let host = Location.query('node');

    if (host === undefined) {
        // if host is null, using default.
        host = hyper.attr('default');
        Location.query('node', host);
    }
    Api.Host(host);

    new Navigation({
        id: "nav",
        home: ".",
        container: container,
        name: hyper.attr('name'),
    });
    new Filters();
    new Routes({
        hyper: hyper,
        container: container,
    });
});

