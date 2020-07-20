import {Container} from "./container.js"
import {Utils} from "../../lib/utils.js";
import {Api} from "../../api/api.js";
import {NetworkApi} from "../../api/network.js";
import {Collapse} from "../collapse.js";
import {NetworkCtl} from "../../ctl/network.js";


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
        console.log('Instance', props);

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
            this.loading();
        });
    }

    loading() {
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

        let net = new NetworkCtl({
            id: this.id(),
            header: {id: this.id("#header")},
            leases: {id: this.id("#leases")},
            subnets: {id: this.id("#subnets")},
        });
        // new InstanceSet({id: '#InstanceSetModal', cpu: instance.cpu, mem: instance.mem })
        //     .onsubmit((e) => {
        //         instance.edit(Utils.toJSON(e.form));
        //     });
    }

    template(v) {
        return this.compile(`
        <div id="network" data="{{uuid}}" name="{{name}}">
        <div id="header" class="card">
            <div class="card-header">
                <div class="card-just-left">
                    <a id="refresh" class="none" href="javascript:void(0)">{{name}}</a>
                </div>
            </div>
            <!-- Overview -->
            <div id="collapseOver" class="collapse" aria-labelledby="headingOne" data-parent="#instance">
            <div class="card-body">
                <!-- Header buttons -->
                <div class="card-body-hdl">
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">{{'refresh' | i}}</button>
                    <button id="autostart" type="button" class="btn btn-outline-dark btn-sm">{{'autostart' | i}}</button>
                    <div id="btns-more" class="btn-group btn-group-sm" role="group">
                        <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle"
                                data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                            {{'actions' | i}}
                        </button>
                        <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                            <a id="edit" class="dropdown-item" href="javascript:void(0)">{{'edit' | i}}</a>
                            <a id="dumpxml" class="dropdown-item" href="javascript:void(0)">{{'dump xml' | i}}</a>
                            <div class="dropdown-divider"></div>
                            <a id="destroy" class="dropdown-item" href="javascript:void(0)">{{'destroy' | i}}</a>
                            <a id="remove" class="dropdown-item" href="javascript:void(0)">{{'remove' | i}}</a>
                        </div>
                    </div>
                </div>
                <dl class="dl-horizontal">
                    <dt>{{'state' | i}}:</dt>
                    <dd><span class="{{state}}">{{state}}</span></dd>
                    <dt>UUID:</dt>
                    <dd>{{uuid}}</dd>
                    <dt>{{'mode' | i}}:</dt>
                    <dd>{{mode == '' ? 'isolated' : mode}}</dd>
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
        <div id="collapse">
        <!-- DHCP Lease -->
        <div id="leases" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseLea" aria-expanded="true" aria-controls="collapseLea">
                    {{'dhcp lease' | i}}
                </button>
            </div>
            <div id="collapseLea" class="collapse" aria-labelledby="headingOne" data-parent="#collapse">
            <div class="card-body">
                <div class="card-body-hdl">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#LeaseCreateModal">
                        {{'new a lease' | i}}
                    </button>
                    <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                    <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
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
        </div>`, v);
    }
}
