import {Container} from "./container.js"
import {CollapseWid} from "../widget/collapse.js";
import {DataStoreApi} from "../api/datastoresapi.js";
import {PoolCtl as PoolController} from "../controller/poolctl.js";
import {Api} from "../api/api.js";
import {VolumeRemoveWid} from "../widget/volume/volumeremove.js";
import {VolumeApi} from "../api/volumeapi.js";


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
        // CollapseWid
        $(this.id('#collapseOver')).fadeIn('slow');
        $(this.id('#collapseOver')).collapse('show');
        new CollapseWid({
            pages: [
                {id: this.id('#collapseVol'), name: 'volumes'},
            ],
            default: this.default,
            update: false,
        });

        new PoolController({
            id: this.id(),
            header: {id: this.id("#header")},
            confirm: this.id("#confirmActionModal"),
            volumes: {
                id: this.id("#volumes"),
                upload: "#uploadPoolModal",
                onRemove: (objs) => {
                    new VolumeRemoveWid({
                        id: this.id('#removeVolModal'),
                        name: objs.uuids,
                    }).onsubmit((e) => {
                        new VolumeApi(objs).delete();
                    });
                }
            },
            upload: '#uploadStoreModal',
        });

    }

    template(v) {
        let dumpUrl = Api.path(`/api/datastore/${v.uuid}?format=xml`);

        return this.compile(`
        <div id="datastores" data="{{uuid}}" name="{{name}}" state="{{state}}" autostart="{{autostart}}">
        <div id="header" class="card shadow">
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
                        <button id="upload" type="button" class="btn btn-outline-info btn-sm"
                                 data-toggle="modal" data-target="#uploadStoreModal">{{'upload file' | i}}</button>
                        <button id="autostart" type="button" class="btn btn-outline-dark btn-sm">
                            {{if autostart}}{{'disable autostart' | i}}{{else}}{{'enable autostart' | i}}{{/if}}
                        </button>
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
                <div class="card-body-tbl pl-1 pr-1 pt-2 pb-2">
                    <div class="resource-overview">
                        <div class="dashboard-grid datastore-overview-grid">
                            <div class="dashboard-stat total">
                                <div class="label">{{'name' | i}}</div>
                                <div class="value">{{name}}</div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'state' | i}}</div>
                                <div class="value"><span class="st-{{state}}">{{state}}</span></div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'allocation' | i}}</div>
                                <div class="value">{{allocation | prettyByte}}</div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'capacity' | i}}</div>
                                <div class="value">{{capacity | prettyByte}}</div>
                            </div>
                            <div class="dashboard-stat resource-wide">
                                <div class="label">{{'uuid' | i}}</div>
                                <div class="value resource-code">{{uuid}}</div>
                            </div>
                            <div class="dashboard-stat resource-wide">
                                <div class="label">{{'source' | i}}</div>
                                <div class="value resource-code">{{source}}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- VolumeCtl list-->
        <div id="volumes" class="card shadow">
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
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#uploadPoolModal">{{'upload' | i}}</button>
                        <button id="remove" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#removeVolModal">{{'remove' | i}}</button>
                    </div>
                    <div class="col-auto mt-auto mb-auto">
                        <button id="datastore" class="btn btn-link btn-sm p-0" data="{{uuid}}">{{name}}:</button>
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
            <div id="removeVolModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="uploadStoreModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="uploadPoolModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="confirmActionModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        </div>`, v);
    }
}
