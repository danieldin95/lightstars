import {Base} from "./base.js"
import {Utils} from "../../com/utils.js";
import {Api} from "../../api/api.js";
import {NetworkApi} from "../../api/network.js";
import {Collapse} from "../collapse.js";
import {NetworkCtl} from "../../ctl/network.js";


export class Network extends Base {
    // {
    //    parent: "#Container",
    //    uuid: "",
    //    default: "lease"
    //    force: false, // force to apply default.
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
            this.view.find(this.id('#refresh')).on('click', (e) => {
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
            leases: {id: this.id("#leases")},
            subnets: {id: this.id("#subnets")},
        });
        // new InstanceSet({id: '#InstanceSetModal', cpu: instance.cpu, mem: instance.mem })
        //     .onsubmit((e) => {
        //         instance.edit(Utils.toJSON(e.form));
        //     });
    }

    template(v) {
        return template.compile(`
        <div id="network" class="card instance" data="{{uuid}}" name="{{name}}">
            <div class="card-header">
                <div class="card-just-left">
                    <a id="refresh" class="none">{{name}}</a>
                </div>
            </div>
            <!-- Overview -->
            <div id="collapseOver" class="collapse" aria-labelledby="headingOne" data-parent="#instance">
            <div class="card-body">
                <!-- Header buttons -->
                <div class="card-header-cnt">
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">Refresh</button>
                    <button id="autostart" type="button" class="btn btn-outline-dark btn-sm">Autostart</button>
                    <div id="btns-more" class="btn-group btn-group-sm" role="group">
                        <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle"
                                data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                            Actions
                        </button>
                        <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                            <a id="edit" class="dropdown-item" href="javascript:void(0)">Edit</a>
                            <a id="destroy" class="dropdown-item" href="javascript:void(0)">Destroy</a>
                            <div class="dropdown-divider"></div>
                            <a id="remove" class="dropdown-item" href="javascript:void(0)">Remove</a>
                            <div class="dropdown-divider"></div>
                            <a id="dumpxml" class="dropdown-item" href="javascript:void(0)">Dump XML</a>
                        </div>
                    </div>
                </div>
                <dl class="dl-horizontal">
                    <dt>State:</dt>
                    <dd><span class="{{state}}">{{state}}</span></dd>
                    <dt>UUID:</dt>
                    <dd>{{uuid}}</dd>
                    <dt>Mode:</dt>
                    <dd>{{mode == '' ? 'isolated' : mode}}</dd>
                    <dt>Address:</dt>
                    <dd>{{if address == ''}} 
                      - 
                    {{else}} 
                      {{address}}/{{if prefix}} {{prefix}} {{else}} {{netmask | netmask2prefix}} {{/if}}
                    {{/if}}</dd>
                </dl>
            </div>
            </div>
        </div>
        <div id="devices">
        <!-- DHCP Lease -->
        <div id="leases" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseLea" aria-expanded="true" aria-controls="collapseLea">
                    DHCP Lease
                </button>
            </div>
            <div id="collapseLea" class="collapse" aria-labelledby="headingOne" data-parent="#devices">
            <div class="card-body">
                <div class="card-header-cnt">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#LeaseCreateModal">
                        New a lease
                    </button>
                    <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                    <button id="remove" type="button" class="btn btn-outline-dark btn-sm">Remove</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                </div>
                <div class="">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th><input id="on-all" type="checkbox"></th>
                            <th>ID</th>
                            <th>MAC address</th>
                            <th>IP address</th>
                            <th>Host name</th>
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
        </div>`)(v);
    }
}
