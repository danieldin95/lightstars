import {Api} from "../api/api.js";
import {HyperApi} from "../api/hyper.js";
import {ZoneApi} from "../api/zone.js";
import {Location} from "../lib/location.js";
import {ChangePassword} from "./user/password/change.js";
import {PasswordApi} from "../api/password.js";
import {Utils} from "../lib/utils.js";
import {Widget} from "./widget.js";
import {Preferences} from "./user/preferences.js";



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
        this.refresh();
    }

    firefox() {
        return navigator.userAgent.match(/firefox/i);
    }

    chrome() {
        return navigator.userAgent.match(/chrome/i);
    }

    nassh(host) {
        let name = 'pnhechapfaindjhompbnflcldabbghjo';
        let extension = `chrome-extension://${name}/html/nassh.html`;
        let webstore = `https://chrome.google.com/webstore/detail/secure-shell-app/${name}`;
        $.get(extension, (resp) => {
            window.open(extension + '#root@' + host);
        }).fail((e)=> {
            window.open(webstore);
        });
    }

    page(url) {
        if (url.indexOf('?') >= 0) {
            return url.split("?", 2)[0];
        }
        return url
    }

    refresh() {
        let page = Location.get("/system");
        this.active = "#" + page;
        console.log("Navigation.refresh", this.active);

        let forceActive = (cur) => {
            this.active = cur;
            console.log("foreActive", cur);
            this.view.find('li').each((i, e) => {
                let href = $(e).find('a').attr("href");
                if (cur && cur === this.page(href)) {
                    $(e).addClass("active");
                } else {
                    $(e).removeClass("active");
                }
            });
        };

        new HyperApi().get(this, (e) => {
            this.view = $(this.render(e.resp));
            forceActive(this.active);
            this.node();
            $(this.parent).html(this.view);

            // register Listener
            this.view.find("#fullscreen").on('click', (e) => {
                this.fullscreen();
            });
            if (this.chrome()) {
                let host = $(location).attr("hostname");
                this.view.find("#phy #ssh").on('click', (e) => {
                    this.nassh(host);
                });
            } else {
                let host = $(location).attr("hostname");
                this.view.find("#phy #ssh").on('click', (e) => {
                    window.open("ssh://"+host);
                });
            }
            let user = e.resp.user.name;
            new ChangePassword({id: '#changePasswdModal'})
                .onsubmit((e) => {
                    new PasswordApi({uuids: user}).set(Utils.toJSON(e.form));
                });
            new Preferences({id: '#preferencesModal'})
                .onsubmit((e) => {
                    console.log(Utils.toJSON(e.form));
                });
        });
    }

    nodeName(host, view) {
        view = view || this.view;
        if (host === "") {
            view.find("node-name").text('localhost');
        } else {
            view.find("node-name").text(host);
        }
    }

    node() {
        let view = this.view;
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
                // reset host.
                Location.query('node', host);
                Api.host(host);
                // fore to instances
                Location.set("/instances");
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
        <nav id="navs" class="navbar fixed-top navbar-expand-md navbar-dark bg-dark">
        <!-- Brand -->
        <a class="navbar-brand" href="${this.home}">
            <img src="/static/images/lightstar.png" width="30" height="30" alt="">
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
                    <a id="system" class="nav-link" href="${this.url('#/system')}">{{'home' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="instances" class="nav-link" href="${this.url('#/instances')}">{{'guest instances' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="datastore" class="nav-link" href="${this.url('#/datastores')}">{{'datastore' | i}}</a>
                </li>
                <li class="nav-item">
                    <a id="network" class="nav-link" href="${this.url('#/networks')}">{{'network' | i}}</a>
                </li>
                <li class="nav-item dropdown">
                    <!-- Physical -->
                    <a id="phyMore" class="nav-link dropdown-toggle" href="javascript:void(0)" 
                        data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        {{'host' | i}} 
                    </a>
                    <div id="phy" class="dropdown-menu" aria-labelledby="phyMore">
                        <a id="disk" class="dropdown-item" href="javascript:void(0);">{{'hardware disk' | i}}</a>
                        <a id="network" class="dropdown-item" href="javascript:void(0);">{{'network interface' | i}}</a>
                        <a id="rootfs" class="dropdown-item" href="javascript:void(0);">{{'filesystem' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a id="ssh" class="dropdown-item" href="javascript:void(0);">{{'ssh console' | i}}</a>
                    </div>
                </li>
            </ul>
            <!-- Zone -->
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="nodeMore" class="nav-link dropdown-toggle" href="javascript:void(0)" data-toggle="dropdown" aria-haspopup="true"
                       aria-expanded="false">
                        <node-name>localhost</node-name>@{{'node' | i}}
                    </a>
                    <div id="node" class="dropdown-menu" aria-labelledby="nodeMore">
                        <a class="dropdown-item" data="">default</a>
                    </div>
                </li>            
            </ul>      
            <!-- User -->      
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="userMore" class="nav-link dropdown-toggle" href="javascript:void(0)" data-toggle="dropdown" 
                        aria-haspopup="true" aria-expanded="false">
                        {{user.name}}@{{hyper.host}}
                    </a>
                    <div class="dropdown-menu dropdown-left" aria-labelledby="userMore">
                        <a id="fullscreen" class="dropdown-item">{{'full screen' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="javascript:void(0)" data-toggle="modal" 
                            data-target="#preferencesModal">{{'preferences' | i}}</a>
                        <a class="dropdown-item" href="javascript:void(0)" 
                            data-toggle="modal" data-target="#changePasswdModal">{{'change password' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="/ui/login">{{'logout' | i}}</a>
                    </div>
                </li>
            </ul>
        </div>
        </nav>
        <!-- Modals -->
        <div id="modals" class="modals">
            <div id="preferencesModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="changePasswdModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        `, v)
    }
}
