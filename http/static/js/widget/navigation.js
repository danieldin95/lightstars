import {Api} from "../api/api.js";
import {HyperApi} from "../api/hyper.js";
import {ZoneApi} from "../api/zone.js";
import {Location} from "../com/location.js";
import {Index} from "../widget/container/index.js";
import {ChangePassword} from "./password/change.js";
import {PasswordApi} from "../api/password.js";
import {Utils} from "../com/utils.js";



export class Navigation {
    // {
    //   id: '#xx'.
    //   home: '.'
    // }
    constructor(props) {
        this.id = props.id;
        this.home = props.home;
        this.props = props;
        this.active = Location.get("instances");
        console.log("Navigation.constructor", this.active);
        this.navs = ["#system", "#instances", "#datastore", "#network"];

        this.refresh();
    }

    refresh() {
        let active = (cur) => {
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
              active(e.name);
            },
        });
        new HyperApi().get(this, (e) => {
            let view = $(this.render(e.resp));
            let name = this.props.name;
            let container = this.props.container;

            this.view = view;
            for (let nav of this.navs) {
                this.view.find(nav).on('click', function (e) {
                    active($(this).parent('li a').attr("id"));
                    new Index({
                        id: container,
                        name: name,
                        force: true,
                        default: $(this).attr("data-target"),
                    });
                })
            }

            active(this.active);
            this.zone();

            this.view.find("#fullscreen").on('click', (e) => {
                this.fullscreen();
            });

            let password = new PasswordApi();
            new ChangePassword({id: '#changePasswordModal'}).onsubmit((e) => {
                password.set(Utils.toJSON(e.form));
            });
            $(this.id).html(this.view);
        });
    }

    zoneName(host, view) {
        view = view || this.view;

        if (host === "") {
            view.find("zone-name").text('default');
        } else {
            view.find("zone-name").text(host);
        }
    }

    zone() {
        let view = this.view;
        let name = this.props.name;
        let container = this.props.container;
        let host = Location.query('node');

        this.zoneName(host);
        new ZoneApi().list(this, (data) => {
            view.find("#zone").empty();
            data.resp.forEach((v, i) => {
                let name = v['name'];
                let value = v['name'];
                if (v['url'] === '') {
                    value = '';
                }
                let elem = $(`
                   <div class="dropdown-divider"></div>
                   <a class="dropdown-item" data="${value}">${name}</a>
                `);
                if (i === 0) {
                    elem = $(`<a class="dropdown-item" data="${value}">${name}</a>`);
                }
                view.find("#zone").append(elem);
            });
            view.find("#zone a").on('click',  this,function (e) {
                let host = $(this).attr("data");

                e.data.zoneName(host, view);
                Location.query('node', host);
                Api.Host(host);
                new Index({
                    id: container,
                    name: name,
                    force: true,
                    default: e.data.active,
                });
            });
        });
    }

    fullscreen() {
        let el = document.documentElement
            , rfs =
            el.requestFullScreen
            || el.webkitRequestFullScreen
            || el.mozRequestFullScreen
            || el.msRequestFullScreen
        ;
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

    render(data) {
        return template.compile(`
        <a class="navbar-brand" href="${this.home}">
            <img src="/static/images/lightstar-6.png" width="30" height="30" alt="">
        </a>
        <button class="navbar-toggler" type="button" data-toggle="collapse"
                data-target="#navbarMore" aria-controls="navbarMore" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"/>
        </button>
        <div class="collapse navbar-collapse" id="navbarMore">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <a id="system" class="nav-link" data-target="/system">Home</a>
                </li>
                <li class="nav-item">
                    <a id="instances" class="nav-link" data-target="/instances">Guest Instances</a>
                </li>
                <li class="nav-item">
                    <a id="datastore" class="nav-link" data-target="/datastore">DataStore</a>
                </li>
                <li class="nav-item">
                    <a id="network" class="nav-link" data-target="/network">Network</a>
                </li>
            </ul>
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="zoneMore" class="nav-link dropdown-toggle" href="#" data-toggle="dropdown" aria-haspopup="true"
                       aria-expanded="false">
                        <zone-name>default</zone-name>@zone
                    </a>
                    <div id="zone" class="dropdown-menu" aria-labelledby="zoneMore">
                        <a class="dropdown-item" data="">default</a>
                    </div>
                </li>            
            </ul>            
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a id="userMore" class="nav-link dropdown-toggle" href="#" data-toggle="dropdown" 
                        aria-haspopup="true" aria-expanded="false">
                        {{user.name}}@{{hyper.host}}
                    </a>
                    <div class="dropdown-menu dropdown-left" aria-labelledby="userMore">
                        <a id="fullscreen" class="dropdown-item">Full screen</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="#">Setting</a>
                        <a class="dropdown-item" href="#" data-toggle="modal" data-target="#changePasswordModal">Change password</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="/ui/login">Logout</a>
                    </div>
                </li>
            </ul>
        </div>`)(data)
    }
}
