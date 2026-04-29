import {Container} from "./container.js"
import {Api} from "../api/api.js";
import {NetworkApi} from "../api/networkapi.js";
import {collapse} from "../widget/collapse.js";
import {Network as NetworkController} from "../controller/networkctl.js";

export class Network extends Container {
    // {
    //    parent: "#container",
    //    uuid: "",
    //    default: "lease"
    // }
    constructor(props) {
        super(props);
        this.default = props.default || 'lease';
        this.uuid = props.uuid;
        this.current = "#network";

        this.render();
    }

    render() {
        new NetworkApi({uuids: this.uuid}).get(this, (e) => {
            this.title(e.resp.name);
            this.view = $(this.template(e.resp));
            this.view.find('#header #refresh').on('click', (e) => {
                this.render();
            });
            $(this.parent).html(this.view);
            this.loading(e.resp);
        });
    }

    loading(data) {
        // collapse
        $(this.id('#collapseOver')).fadeIn('slow');
        $(this.id('#collapseOver')).collapse();
        new collapse({
            pages: [
                {id: this.id('#collapseLea'), name: 'lease'},
            ],
            default: this.default,
            update: false,
        });

        new NetworkController({
            id: this.id(),
            header: {id: this.id("#header")},
            confirm: this.id("#confirmActionModal"),
            leases: {id: this.id("#leases")},
            ports: {
                id: this.id("#ports"),
                bridge: data.bridge,
            }
        });
    }

    template(v) {
        let dumpUrl = Api.path(`/api/network/${v.uuid}?format=xml`);

        return this.compile(`
        <div id="network" data="{{uuid}}" name="{{name}}" state="{{state}}" autostart="{{autostart}}">
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
                <div class="card-body-hdl pt-1">
                    <div class="resource-overview">
                        <div class="dashboard-grid network-overview-grid">
                            <div class="dashboard-stat total">
                                <div class="label">{{'name' | i}}</div>
                                <div class="value">{{name}}</div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'state' | i}}</div>
                                <div class="value"><span class="st-{{state}}">{{state}}</span></div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'mode' | i}}</div>
                                <div class="value">{{mode == '' ? 'isolated' : mode}}</div>
                            </div>
                            <div class="dashboard-stat total">
                                <div class="label">{{'bridge' | i}}</div>
                                <div class="value">{{bridge}}</div>
                            </div>
                            <div class="dashboard-stat resource-wide">
                                <div class="label">UUID</div>
                                <div class="value resource-code">{{uuid}}</div>
                            </div>
                            <div class="dashboard-stat resource-wide">
                                <div class="label">{{'address' | i}}</div>
                                <div class="value resource-code">
                                    {{if address == ''}}
                                        -
                                    {{else}}
                                        {{address}}/{{if prefix}}{{prefix}}{{else}}{{netmask | netmask2prefix}}{{/if}}
                                    {{/if}}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="card-tab">
            <ul class="nav nav-pills justify-content-start" id="pills-tab" role="tablist">
              <li class="nav-item" role="presentation">
                <a class="nav-link active" id="pills-0-tab" data-toggle="pill" href="#pills-0" 
                    role="tab" aria-controls="pills-0" aria-selected="true">{{'virtual ports' | i}}</a>
              </li>
              <li class="nav-item" role="presentation">
                <a class="nav-link" id="pills-1-tab" data-toggle="pill" href="#pills-1" 
                    role="tab" aria-controls="pills-1" aria-selected="false">{{'dhcp lease' | i}}</a>
              </li>
            </ul>
            <div class="tab-content" id="pills-tabContent">
              <div class="tab-pane fade show active" id="pills-0" role="tabpanel" aria-labelledby="pills-0-tab">
                <!-- virtual Ports -->
                <div id="ports" class="card shadow">
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                        data-toggle="modal" data-target="#PortCreateModal">
                                    {{'create port' | i}}
                                </button>
                                <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                                <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
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
                                    <th>{{'id' | i}}</th>
                                    <th>{{'instance' | i}}</th>
                                    <th>{{'device' | i}}</th>
                                    <th>{{'mac' | i}}</th>
                                    <th>{{'ip address' | i}}</th>
                                    <th>{{'model' | i}}</th>
                                </tr>
                                </thead>
                                <tbody id="display-table">
                                <!-- Loading... -->
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>             
              </div>
              <div class="tab-pane fade" id="pills-1" role="tabpanel" aria-labelledby="pills-1-tab">
                <!-- DHCP Lease -->
                <div id="leases" class="card shadow">
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                        data-toggle="modal" data-target="#LeaseCreateModal">
                                    {{'new a lease' | i}}
                                </button>
                                <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                                <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
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
                                    <th>{{'id' | i}}</th>
                                    <th>{{'hostname' | i}}</th>
                                    <th>{{'mac' | i}}</th>
                                    <th>{{'ip address' | i}}</th>
                                </tr>
                                </thead>
                                <tbody id="display-table">
                                <!-- Loading... -->
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
              </div>
            </div>
        </div>
        <div id="modals">
            <div id="confirmActionModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        </div>`, v);
    }
}
