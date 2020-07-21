import {Container} from "./container.js"
import {Location} from "../../lib/location.js";
import {Overview} from "../index/overview.js";
import {I18N} from "../../lib/i18n.js";

export class Home extends Container {
    // {
    //    parent: "#container",
    // }
    constructor(props) {
        super(props);
        this.current = '#index';

        this.render();
        this.loading();
    }

    loading() {
        this.title(I18N.i('home'));
        // loading overview.
        let view = new Overview({
            id: this.id('#overview'),
        });
        view.refresh((e) => {
            this.props.name = e.resp.hyper.name;
            $(this.id('#refresh')).text(this.props.name);
        });
        // register click on overview.
        $(this.id('#refresh')).on('click', () => {
            view.refresh();
        });
    }

    template(v) {
        return this.compile(`
        <div id="index">
        <!-- System -->
        <div id="system" class="card card-main system">
            <div class="card-header">
                <button id="refresh" class="btn btn-link btn-block text-left btn-sm">${this.props.name}</button>
            </div>
            <div id="collapseSys">
            <div id="overview" class="card-body">
            <!-- Loading -->
            </div>
            </div>
        </div>
        
        <!-- Modal -->
        <div id="modals">
        </div>
        </div>`)
    }
}
