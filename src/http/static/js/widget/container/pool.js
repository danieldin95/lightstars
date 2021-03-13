import {Container} from "./container.js"
import {Collapse} from "../collapse.js";
import {DataStoreApi} from "../../api/datastores.js";
import {PoolCtl} from "../../ctl/pool.js";
import {Api} from "../../api/api.js";


export class Pool extends Container {

    constructor(props) {
        super(props);
        this.default = props.default || 'volumes';
        this.uuid = props.uuid;
        this.current = "#datastores";

        this.render();
    }

    render() {
        new DataStoreApi({uuids: this.uuid}).get(this, (e) => {
            this.title(e.resp.name);
            this.view = $(this.template(e.resp));
            this.view.find('#header #refresh').on('click', (e) => {
                this.render();
            });
            $(this.parent).html(this.view);
            this.loading();
        });
    }

    loading() {
        // collapse
        $(this.id('#collapseOver')).fadeIn('slow');
        $(this.id('#collapseOver')).collapse();
        new Collapse({
            pages: [
                {id: this.id('#collapseVol'), name: 'volumes'},
            ],
            default: this.default,
            update: false,
        });

        new PoolCtl({
            id: this.id(),
            header: {id: this.id("#header")},
            volumes: {
                id: this.id("#volumes"),
                upload: "#uploadPoolModal",
            },
            upload: '#uploadStoreModal',
        });

    }

    template(v) {
        let dumpUrl = Api.path(`/api/datastore/${v.uuid}?format=xml`);

        return this.compile(`
        <div id="datastores" data="{{uuid}}" name="{{name}}">
        <div id="header" class="card">
            <div class="card-header">
                <div class="text-left">
                    <a id="refresh" class="none" href="javascript:void(0)">{{name}}</a>
                </div>
            </div>
            <!-- Overview -->
            <div class="card-body">
                <!-- Header buttons -->
                <div class="row card-body-hdl">
                    <div class="col-auto mr-auto">
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm" 
                                 data-toggle="modal" data-target="#uploadStoreModal">{{'upload' | i}}</button>
                        <div id="btns-more" class="btn-group btn-group-sm" role="group">
                            <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle"
                                    data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                                {{'actions' | i}}
                            </button>
                            <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                                <a id="edit" class="dropdown-item" href="javascript:void(0)">{{'edit' | i}}</a>
                                <a id="dumpxml" class="dropdown-item" href="${dumpUrl}">{{'dump xml' | i}}</a>
                                <div class="dropdown-divider"></div>
                                <a id="destroy" class="dropdown-item" href="javascript:void(0)">{{'destroy' | i}}</a>
                                <a id="remove" class="dropdown-item" href="javascript:void(0)">{{'remove' | i}}</a>
                            </div>
                        </div>
                    </div>
                    <div class="col-auto">
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">{{'refresh' | i}}</button>                    
                    </div>
                </div>
                <div class="card-body-hdl">
                    <div class="overview">
                        <dl class="dl-horizontal">
                            <dt>{{'state' | i}}:</dt>
                            <dd><span class="{{state}}">{{state}}</span></dd>
                            <dt>{{'uuid' | i}}:</dt>
                            <dd>{{uuid}}</dd>
                            <dt>{{'source' | i}}:</dt>
                            <dd>{{source}}</dd>
                            <dt>{{'allocation' | i}}:</dt>
                            <dd>{{allocation | prettyByte}}</dd>
                            <dt>{{'capacity' | i}}:</dt>
                            <dd>{{capacity | prettyByte}}</dd>
                        </dl>
                    </div>
                </div>
            </div>
        </div>
        <!-- Volume list-->
        <div id="volumes" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseVol" aria-expanded="true" aria-controls="collapseVol">
                    {{'file browser' | i}}
                </button>
            </div>
            <!-- volume actions button-->
            <div class="card-body">
                <div class="row card-body-hdl">
                    <div class="col-auto mr-auto">
                        <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#createVolumeModal">
                            {{'new a volume' | i}}
                        </button>
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm" 
                                data-toggle="modal" data-target="#uploadPoolModal">{{'upload' | i}}</button>
                        <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
                    </div>
                    <div class="col-auto">
                        <button id="datastore" class="btn btn-link btn-sm p-0" data="{{uuid}}">{{name}}:/</button>
                        <button id="current"  class="btn btn-link btn-sm p-0 pr-2" data=""></button>
                    </div>
                    <div class="col-auto">
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    </div>
                </div>
                <div class="card-body-tbl">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th><input id="on-all" type="checkbox"></th>
                            <th>{{'type' | i}}</th>
                            <th>{{'name' | i}}</th>
                            <th>{{'capacity' | i}}</th>
                            <th>{{'allocation' | i}}</th>
                        </tr>
                        </thead>
                        <tbody id="display-table">
                        <!-- Loading... -->
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <div id="modals">
            <div id="uploadStoreModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="uploadPoolModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        </div>`, v);
    }
}
