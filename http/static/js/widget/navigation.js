import {HyperApi} from "../api/hyper.js";


export class Navigation {
    // {
    //   id: '#xx'.
    //   home: '.'
    // }
    constructor(props) {
        this.id = props.id;
        this.home = props.home;
        this.props = props;
        this.active = this.get(window.location.href, "xx");

        this.refresh();
    }

    refresh() {
        new HyperApi().get(this, (e) => {
            this.view = $(this.render(e.resp));
            this.view.find('li').each((i, e) => {
               let href = $(e).find('a').attr("href");
               if (this.get(href) == this.active) {
                   $(e).addClass("active");
               }
            });
            this.view.find("#fullscreen").on('click', (e) => {
                this.fullscreen();
            });
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
        if (typeof rfs != "undefined" && rfs) {
            rfs.call(el);
        } else if (typeof window.ActiveXObject != "undefined") {
            // for Internet Explorer
            let wscript = new ActiveXObject("WScript.Shell");
            if (wscript != null) {
                wscript.SendKeys("{F11}");
            }
        }
    }

    get (href, name) {
        let path = href.split("?", 2)[0];
        let pages = path.split('#', 2);

        name = name || ""
        return (pages.length == 2 && pages[1] != "") ? pages[1] : name;
    }

    render(data) {
        return template.compile(`
        <a class="navbar-brand" href="${this.home}">
            <img src="/static/images/lightstar-6.png" width="30" height="30" alt="">
        </a>
        <button class="navbar-toggler" type="button" data-toggle="collapse"
                data-target="#navbarMore" aria-controls="navbarMore" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarMore">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <a class="nav-link" href="${this.home}#system">Home</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="${this.home}#instances">Guest Instances</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="${this.home}#datastore">DataStore</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="${this.home}#network">Network</a>
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