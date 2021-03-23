import {Container} from "./container.js"
import {Api} from "../api/api.js";
import {NetworkApi} from "../api/network.js";
import {Collapse} from "../widget/collapse.js";
import {NetworkCtl} from "../controller/network.js";

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
        new Collapse({
            pages: [
                {id: this.id('#collapseLea'), name: 'lease'},
            ],
            default: this.default,
            update: false,
        });

        new NetworkCtl({
            id: this.id(),
            header: {id: this.id("#header")},
            leases: {id: this.id("#leases")},
            ports: {
                id: this.id("#ports"),
                uuid: data.bridge,
            }
        });
    }

    template(v) {
        let dumpUrl = Api.path(`/api/network/${v.uuid}?format=xml`);

        return this.compile(`
        <div id="network" data="{{uuid}}" name="{{name}}">
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
                        <button id="autostart" type="button" class="btn btn-outline-dark btn-sm">{{'autostart' | i}}</button>
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
                            <dt>{{'name' | i}}:</dt>
                            <dd>{{name}}</dd>                            
                            <dt>{{'state' | i}}:</dt>
                            <dd><span class="st-{{state}}">{{state}}</span></dd>
                            <dt>UUID:</dt>
                            <dd>{{uuid}}</dd>
                            <dt>{{'mode' | i}}:</dt>
                            <dd>{{mode == '' ? 'isolated' : mode}}</dd>
                            <dt>{{'bridge' | i}}:</dt>
                            <dd>{{bridge}}</dd>                            
                            <dt>{{'address' | i}}:</dt>
                            <dd>{{if address == ''}} 
                              - 
                            {{else}} 
                              {{address}}/{{if prefix}} {{prefix}} {{else}} {{netmask | netmask2prefix}} {{/if}}
                            {{/if}}</dd>
                        </dl>
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
                                <button id="create" type="button" class="btn btn-outline-success btn-sm"
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
                                <button id="create" type="button" class="btn btn-outline-success btn-sm"
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
        </div>`, v);
    }
}
