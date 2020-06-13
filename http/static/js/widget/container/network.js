import {Base} from "./base.js"
import {Utils} from "../../com/utils.js";
import {Api} from "../../api/api.js";
import {NetworkApi} from "../../api/network.js";
import {Collapse} from "../collapse.js";


export class Network extends Base {
    // {
    //    id: ".container",
    //    uuid: "",
    //    default: "lease"
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        super(props);
        this.default = props.default || 'lease';
        this.uuid = props.uuid;
        console.log('Instance', props);

        this.render();
    }

    render() {
        new NetworkApi({uuids: this.uuid}).get(this, (e) => {
            this.title(`${e.resp.name} - LightStar`);
            this.view = $(this.template(e.resp));
            this.view.find('#network #refresh').on('click', (e) => {
                this.render();
            });
            $(this.id).html(this.view);
            this.loading();
        });
    }

    loading() {
        // collapse
        $('#collapseOver').fadeIn('slow');
        $('#collapseOver').collapse();
        new Collapse({
            pages: [
                {id: '#collapseLea', name: 'lease'},
            ],
            default: this.default,
            update: false,
        });

        // let instance = new GuestCtl({
        //     id: "#instance",
        //     disks: {id: "#disk"},
        //     interfaces: {id: "#interface"},
        //     graphics: {id: "#graphics"},
        // });
        // new InstanceSet({id: '#instanceSetModal', cpu: instance.cpu, mem: instance.mem })
        //     .onsubmit((e) => {
        //         instance.edit(Utils.toJSON(e.form));
        //     });
        // loading lease.
        // new DiskCreate({id: '#diskCreateModal'})
        //     .onsubmit((e) => {
        //         instance.disk.create(Utils.toJSON(e.form));
        //     });
    }

    template(v) {
        let cls = "enable";

        return template.compile(`
        <network>
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
                    <div id="power-btns" class="btn-group btn-group-sm" role="group">
                        <button id="start" type="button" class="btn btn-outline-dark">
                            Power on
                        </button>
                        <button id="power" type="button"
                                class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                data-toggle="dropdown" aria-expanded="false">
                            <span class="sr-only">Toggle Dropdown</span>
                        </button>
                        <div id="power-more" class="dropdown-menu" aria-labelledby="power">
                            <a id="start" class="dropdown-item" href="javascript:void(0)">Power on</a>
                            <a id="shutdown" class="dropdown-item" href="javascript:void(0)">Power off</a>
                            <div class="dropdown-divider"></div>
                            <a id="reset" class="dropdown-item" href="javascript:void(0)">Reset</a>
                            <div class="dropdown-divider"></div>
                            <a id="destroy" class="dropdown-item" href="javascript:void(0)">Destroy</a>
                        </div>
                    </div>
                    <div id="btns-more" class="btn-group btn-group-sm" role="group">
                        <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle"
                                data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                            Actions
                        </button>
                        <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                            <a id="suspend" class="dropdown-item ${cls}" href="javascript:void(0)">Suspend</a>
                        </div>
                    </div>
                </div>
                <dl class="dl-horizontal">
                    <dt>State:</dt>
                    <dd><span class="{{state}}">{{state}}</span></dd>
                    <dt>UUID:</dt>
                    <dd>{{uuid}}</dd>
                    <dt>Arch:</dt>
                    <dd>{{arch}} | {{type}}</dd>
                    <dt>Processor:</dt>
                    <dd>{{maxCpu}} | {{cpuTime}}ms</dd>
                    <dt>Memory:</dt>
                    <dd>{{maxMem | prettyKiB}} | {{memory | prettyKiB}}</dd>
                </dl>
            </div>
            </div>
        </div>
        <div id="devices">
        <!-- DHCP Lease -->
        <div id="lease" class="card device">
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
                            data-toggle="modal" data-target="#diskCreateModal">
                        Attach disk
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
                            <th>Bus</th>
                            <th>Device</th>
                            <th>Source</th>
                            <th>Capacity</th>
                            <th>Available</th>
                            <th>Address</th>
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
        </network>`)(v);
    }
}
