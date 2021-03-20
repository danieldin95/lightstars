import {Api} from "./api/api.js";
import {Location} from "./lib/location.js";
import {I18N} from "./lib/i18n.js";
import {Navigation} from "./widget/navigation.js";
import {Container} from "./container/container.js";
import {Routes} from "./routes.js";
import {Template} from "./lib/template.js";


$(function() {
    let hyper = $('hyper');
    let alias = hyper.attr('alias');
    let host = Location.query('node');
    if (!host) {
        // if host is null, using default.
        host = hyper.attr('default');
        Location.query('node', host);
    }
    Api.host(host);
    Container.alias(alias);

    I18N.promise().then(function () {
        let tmpl = new Template();
        let nav = new Navigation({
            parent: "#navigation",
            home: ".",
            container: "#container",
            hyper: hyper,
        });
        let rte = new Routes({
            hyper: hyper,
            container: "#container",
            onchange: function (e) {
                // remove backdrop of modal.
                $('.modal-backdrop').remove();
                // refresh navigation
                nav.refresh();
            },
        });
    });
});

