import {Api} from "./api/api.js";
import {Location} from "./lib/location.js";
import {I18N} from "./lib/i18n.js";
import {NavigationWid} from "./widget/navigation.js";
import {Container} from "./container/container.js";
import {Routes} from "./routes.js";
import {Template} from "./lib/template.js";
import {Alert} from "./lib/alert.js";
import {HistoryWid} from "./widget/index/history.js";

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
        let renderGlobalHistory = () => {
            $("#global-history").html(`
                <div id="history" class="card shadow history">
                    <div class="card-header">
                        <button class="btn btn-link btn-block text-left btn-sm">${$.i18n("operation history")}</button>
                    </div>
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <input id="search-input" type="text" class="form-control form-control-sm" autocomplete="off" placeholder="${$.i18n("search")}">
                                <button id="download" type="button" class="btn btn-outline-dark btn-sm">${$.i18n("download")}</button>
                            </div>
                            <div class="col-auto">
                                <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">${$.i18n("refresh")}</button>
                            </div>
                        </div>
                        <div class="card-body-tbl">
                            <table class="table table-striped">
                                <thead>
                                <tr>
                                    <th>${$.i18n("user")}</th>
                                    <th>${$.i18n("date")}</th>
                                    <th>${$.i18n("client")}</th>
                                    <th>${$.i18n("message")}</th>
                                    <th>${$.i18n("url")}</th>
                                </tr>
                                </thead>
                                <tbody id="display-table"></tbody>
                            </table>
                        </div>
                        <div class="row mt-2">
                            <div class="col-auto ml-auto d-flex align-items-center">
                                <button id="page-prev" type="button" class="btn btn-outline-dark btn-sm mr-2">${$.i18n("previous")}</button>
                                <span id="page-info" class="smaller mr-2">1/1</span>
                                <button id="page-next" type="button" class="btn btn-outline-dark btn-sm">${$.i18n("next")}</button>
                            </div>
                        </div>
                    </div>
                </div>
            `);
            let his = new HistoryWid({
                id: "#global-history #history .card-body-tbl #display-table",
                pageSize: 10,
                pager: {
                    prev: "#global-history #history #page-prev",
                    next: "#global-history #history #page-next",
                    info: "#global-history #history #page-info",
                },
            });
            his.refresh();
            $("#global-history #history #refresh").off("click").on("click", () => his.refresh());
            $("#global-history #history #download").off("click").on("click", () => his.downloadCSV());
            $("#global-history #history #search-input").off("input").on("input", function () {
                his.setSearch($(this).val() || "");
            });
        };
        renderGlobalHistory();

        let nav = new NavigationWid({
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
                if ($("#global-history #history").length === 0) {
                    renderGlobalHistory();
                }
                $("#global-history #history .card-body-tbl #display-table").trigger("history:refresh");
            },
        });
    });
});
