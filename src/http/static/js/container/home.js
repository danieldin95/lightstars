import {Container} from "./container.js"
import {Overview} from "../widget/index/overview.js";
import {I18N} from "../lib/i18n.js";
import {InstanceCreate} from "../widget/instance/create.js";
import {Utils} from "../lib/utils.js";
import {InstanceApi} from "../api/instance.js";

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
            id: this.id('#overview .card-body-tbl'),
        });
        view.refresh((e) => {
            this.props.name = e.resp.hyper.name;
            $(this.id('#refresh-hdl')).text(this.props.name);
        });
        // register click on overview.
        $(this.id('#refresh')).on('click', () => {
            view.refresh();
        });
        $(this.id('#refresh-hdl')).on('click', () => {
            view.refresh();
        });
        new InstanceCreate({id: '#createGuestModal'})
            .onsubmit((e) => {
                new InstanceApi().create(Utils.toJSON(e.form));
            });
    }

    template(v) {
        return this.compile(`
        <div id="index">
        <!-- System -->
        <div id="system" class="card card-main system">
            <div class="card-header">
                <button id="refresh-hdl" class="btn btn-link btn-block text-left btn-sm">${this.props.name}</button>
            </div>
            <div id="overview" class="card-body">
                <!-- Overview buttons -->
                <div class="row card-body-hdl">
                    <div class="col-auto mr-auto">
                        <button id="create" type="button" class="btn btn-outline-success btn-sm"
                                data-toggle="modal" data-target="#createGuestModal">
                            {{'create new instance' | i}}
                        </button>
                        <button id="console" type="button" class="btn btn-outline-dark btn-sm">{{'power off' | i}}</button>
                        <button id="start" type="button" class="btn btn-outline-dark btn-sm">{{'reboot' | i}}</button>
                    </div>
                    <div class="col-auto">
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    </div>
                </div>
                <div class="card-body-tbl">
                    <!-- Loading -->
                </div>
            </div>
        </div>
        
        <!-- Modal -->
        <div id="modals">
            <!-- Create instance modal -->
            <div id="createGuestModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div> 
        </div>
        </div>`)
    }
}
