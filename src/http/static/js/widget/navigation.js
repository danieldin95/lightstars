import {Api} from "../api/api.js";
import {HyperApi} from "../api/hyper.js";
import {ZoneApi} from "../api/zone.js";
import {Location} from "../com/location.js";
import {Home} from "../widget/container/home.js";
import {ChangePassword} from "./password/change.js";
import {PasswordApi} from "../api/password.js";
import {Utils} from "../com/utils.js";
import {Widget} from "./widget.js";



export class Navigation extends Widget {
    // {
    //   parent: '#xx'.
    //   home: '.'
    // }
    constructor(props) {
        super(props);
        this.parent = props.parent;
        this.home = props.home;
        this.active = "";
        this.navIds = ["#system", "#instances", "#datastore", "#network"];
        this.refresh();
    }

    refresh() {
        this.active = Location.get("/instances");
        console.log("Navigation.refresh", this.active);

        let forceActive = (cur) => {
            this.active = cur;
            this.view.find('li').each((i, e) => {
                let href = $(e).find('a').attr("data-target");
                if (cur && cur === href) {
                    $(e).addClass("active");
                } else {
                    $(e).removeClass("active");
                }
            });
        };

        Location.listen({
            data: this,
            func: (e) => {
                forceActive(e.name);
            },
        });
        new HyperApi().get(this, (e) => {
            let view = $(this.render(e.resp));
            let name = this.props.name;
            let container = this.props.container;

            this.view = view;
            for (let nav of this.navIds) {
                this.view.find(nav).on('click', function (e) {
                    forceActive($(this).parent('li a').attr("id"));
                    new Home({
                        parent: container,
                        name: name,
                        force: true,
                        default: $(this).attr("data-target"),
                    });
                })
            }
            forceActive(this.active);
            this.node();
            $(this.parent).html(this.view);

            // register Listener
            this.view.find("#fullscreen").on('click', (e) => {
                this.fullscreen();
            });
            new ChangePassword({id: '#changePasswdModal'})
                .onsubmit((e) => {
                    new PasswordApi().set(Utils.toJSON(e.form));
                });
        });
    }

    nodeName(host, view) {
        view = view || this.view;

        if (host === "") {
            view.find("node-name").text('default');
        } else {
            view.find("node-name").text(host);
        }
    }

    node() {
        let view = this.view;
        let name = this.props.name;
        let container = this.props.container;
        let host = Location.query('node');

        this.nodeName(host);
        new ZoneApi().list(this, (data) => {
            view.find("#node").empty();
            data.resp.forEach((v, i) => {
                let name = v['name'];
                let value = v['name'];
                if (v['url'] === '') {
                    value = '';
                }
                let elem = $(`
                   <a class="dropdown-item" data="${value}">${name}</a>
                `);
                if (i === 0) {
                    elem = $(`<a class="dropdown-item" data="${value}">${name}</a>`);
                }
                view.find("#node").append(elem);
            });
            view.find("#node a").on('click', this, function (e) {
                let host = $(this).attr("data");

                //e.data.nodeName(host, view);
                Location.set("/instances");
                Location.query('node', host);
                Api.host(host);

                e.data.refresh();
                new Home({
                    parent: container,
                    name: name,
                    force: true,
                    default: e.data.active,
                });
            });
        });
    }

    fullscreen() {
        let el = document.documentElement;
        let rfs = el.requestFullScreen
            || el.webkitRequestFullScreen
            || el.mozRequestFullScreen
            || el.msRequestFullScreen;
        if (rfs) {
            rfs.call(el);
        } else if (window.ActiveXObject) {
            // for Internet Explorer
            let ws = new ActiveXObject("WScript.Shell");
            if (ws != null) {
                ws.SendKeys("{F11}");
            }
        }
    }

    render(v) {
        return this.compile(`
        <nav id="navs" class="navbar sticky-top navbar-expand-lg navbar-dark bg-dark">
        <!-- Brand -->
        <a class="navbar-brand" href="${this.home}">
            <img src="/static/images/lightstar-6.png" width="30" height="30" alt="">
        </a>
        <!-- Collapse bar -->
        <button class="navbar-toggler" type="button" data-toggle="collapse"
                data-target="#navbarMore" aria-controls="navbarMore" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"/>
        </button>
        <!-- More bar -->
        <div class="collapse navbar-collapse" id="navbarMore">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <a id="system" class="nav-link" data-target="/system">{{'home' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="instances" class="nav-link" data-target="/instances">{{'guest instances' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="datastore" class="nav-link" data-target="/datastore">{{'datastore' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="network" class="nav-link" data-target="/network">{{'network' | i}}</a>
                </li>
            </ul>
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="nodeMore" class="nav-link dropdown-toggle" href="javascript:void(0)" data-toggle="dropdown" aria-haspopup="true"
                       aria-expanded="false">
                        <node-name>default</node-name>@{{'node' | i}}
                    </a>
                    <div id="node" class="dropdown-menu" aria-labelledby="nodeMore">
                        <a class="dropdown-item" data="">default</a>
                    </div>
                </li>            
            </ul>            
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="userMore" class="nav-link dropdown-toggle" href="javascript:void(0)" data-toggle="dropdown" 
                        aria-haspopup="true" aria-expanded="false">
                        {{user.name}}@{{hyper.host}}
                    </a>
                    <div class="dropdown-menu dropdown-left" aria-labelledby="userMore">
                        <a id="fullscreen" class="dropdown-item">{{'full screen' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="javascript:void(0)">{{'preferences' | i}}</a>
                        <a class="dropdown-item" href="javascript:void(0)" 
                            data-toggle="modal" data-target="#changePasswdModal">
                            {{'change password' | i}}
                        </a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="/ui/login">{{'logout' | i}}</a>
                    </div>
                </li>
            </ul>
        </div>
        </nav>
        <!-- Modals -->
        <div id="modals" class="modals">
            <div id="changePasswdModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        `, v)
    }
}
