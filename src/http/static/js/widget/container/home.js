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
            $(this.id('#system-col')).text(this.props.name);
        });
        // register click on overview.
        $(this.id('#system-ref')).on('click', () => {
            view.refresh();
        });
    }

    template(v) {
        return this.compile(`
        <div id="index">
        <!-- System -->
        <div id="system" class="card card-main system">
            <div class="card-header">
                <div class="">
                    <a id="system-col" href="javascript:void(0)">${this.props.name}</a>
                    <a class="btn-spot float-right" id="system-ref" href="${this.url('#/system')}"></a>
                </div>
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
