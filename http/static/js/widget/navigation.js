import {HyperApi} from "../api/hyper.js";
import {Location} from "../com/location.js";
import {Index} from "../widget/container/index.js"

export class Navigation {
    // {
    //   id: '#xx'.
    //   home: '.'
    // }
    constructor(props) {
        this.id = props.id;
        this.home = props.home;
        this.props = props;
        this.active = Location.get("xx");
        this.navs = ["#system", "#instances", "#datastore", "#network"];

        this.refresh();
    }


    refresh() {
        let active = (cur) => {
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
            this.view = $(this.render(e.resp));
            active(this.active);

            this.view.find("#fullscreen").on('click', (e) => {
                this.fullscreen();
            });

            let name = this.props.name;
            for (let i = 0; i < this.navs.length; i++) {
                this.view.find(this.navs[i]).on('click', function (e) {
                    active();
                    $(this).parent('li').addClass("active");
                    new Index({
                        id: ".container",
                        name: name,
                        force: true,
                        default: $(this).attr("data-target"),
                    });
                })
            }
            $(this.id).html(this.view);
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
                    <a id="system" class="nav-link" data-target="system">Home</a>
                </li>
                <li class="nav-item">
                    <a id="instances" class="nav-link" data-target="instances">Guest Instances</a>
                </li>
                <li class="nav-item">
                    <a id="datastore" class="nav-link" data-target="datastore">DataStore</a>
                </li>
                <li class="nav-item">
                    <a id="network" class="nav-link" data-target="network">Network</a>
                </li>
            </ul>
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown-toggle" href="#" id="navbarMore" data-toggle="dropdown" aria-haspopup="true"
                       aria-expanded="false">
                        {{user.name}}@{{hyper.host}}
                    </a>
                    <div class="dropdown-menu dropdown-left" aria-labelledby="navbarMore">
                        <a id="fullscreen" class="dropdown-item">Full screen</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="#">Setting</a>
                        <a class="dropdown-item" href="#">Change password</a>
                        <div class="dropdown-divider"></div>
                        <a class="dropdown-item" href="/ui/login">Logout</a>
                    </div>
                </li>
            </ul>
        </div>`)(data)
    }
}