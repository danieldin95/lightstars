import {Base} from "./base.js"
import {Utils} from "../../com/utils.js";
import {Instance} from '../../ctl/instance.js';
import {Api} from "../../api/api.js";
import {InstanceApi} from "../../api/instance.js";
import {Collapse} from "../collapse.js";
import {DiskCreate} from '../disk/create.js';
import {InterfaceCreate} from '../interface/create.js';
import {InstanceSet} from "../instance/setting.js";


export class Guest extends Base {
    // {
    //    id: ".container",
    //    uuid: "",
    //    default: "disk"
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        super(props);
        this.default = props.default || 'disk';
        this.uuid = props.uuid;
        console.log('Instance', props);

        this.render();
    }

    render() {
        new InstanceApi({uuids: this.uuid}).get(this, (e) => {
            this.title(`${e.resp.name} - LightStar`);

            this.view = $(this.template(e.resp));
            this.view.find('#instance #refresh').on('click', (e) => {
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
                {id: '#collapseInt', name: 'interface'},
                {id: '#collapseDis', name: 'disk'},
                {id: '#collapseGra', name: 'graphics'},
            ],
            default: this.default,
            update: false,
        });

        let instance = new Instance({
            id: "#instance",
            disks: {id: "#disk"},
            interfaces: {id: "#interface"},
            graphics: {id: "#graphics"},
        });
        new InstanceSet({id: '#instanceSetModal', cpu: instance.cpu, mem: instance.mem })
            .onsubmit((e) => {
                instance.edit(Utils.toJSON(e.form));
            });

        // loading disks and interfaces.
        new DiskCreate({id: '#diskCreateModal'})
            .onsubmit((e) => {
                instance.disk.create(Utils.toJSON(e.form));
            });
        new InterfaceCreate({id: '#interfaceCreateModal'})
            .onsubmit((e) => {
                instance.interface.create(Utils.toJSON(e.form));
            });

        // register console draggable.
        $(function (e) {
            $('#consoleModal').draggable();
        });
    }

    template(v) {
        let cls = 'disabled';
        let vncUrl = '#';
        let hostUrl = '';

        if (Api.host !== '') {
           hostUrl = "/host/" + Api.host;
        }
        if (v.state === 'running') {
            cls = '';
            let vnc = Utils.graphic(v, 'vnc', 'password');
            vncUrl = "/ui/console?id=" + v.uuid + "&password=" + vnc + "&node=" + Api.host;
        }

        return template.compile(`
        <instance>
        <div id="instance" class="card instance" data="{{uuid}}" name="{{name}}" cpu="{{maxCpu}}" mem="{{maxMem}}">
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
                    <div id="console-btns" class="btn-group btn-group-sm" role="group">
                        <button id="console" type="button" class="btn btn-outline-dark ${cls}"
                                data-target="#consoleModal" data="${vncUrl}">
                            Console
                        </button>
                        <button id="consoles" type="button"
                                class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split ${cls}"
                                data-toggle="dropdown" aria-expanded="false">
                            <span class="sr-only">Toggle Dropdown</span>
                        </button>
                        <div id="console-more" class="dropdown-menu" aria-labelledby="consoles">
                            <a id="console-self" class="dropdown-item" href="#" data="${vncUrl}">
                                Console in self
                            </a>
                            <a id="console-blank" class="dropdown-item" href="#" data="${vncUrl}">
                                Console in new blank
                            </a>
                            <a id="console-window" class="dropdown-item" href="#" data="${vncUrl}">
                                Console in new window
                            </a>
                        </div>
                    </div>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">Refresh</button>
                    <div id="power-btns" class="btn-group btn-group-sm" role="group">
                        <button id="start" type="button" class="btn btn-outline-dark">
                            Power on
                        </button>
                        <button id="power" type="button"
                                class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split ${cls}"
                                data-toggle="dropdown" aria-expanded="false">
                            <span class="sr-only">Toggle Dropdown</span>
                        </button>
                        <div id="power-more" class="dropdown-menu" aria-labelledby="power">
                            <a id="start" class="dropdown-item" href="#">Power on</a>
                            <a id="shutdown" class="dropdown-item" href="#">Power off</a>
                            <div class="dropdown-divider"></div>
                            <a id="reset" class="dropdown-item" href="#">Reset</a>
                            <div class="dropdown-divider"></div>
                            <a id="destroy" class="dropdown-item" href="#">Destroy</a>
                        </div>
                    </div>
                    <div id="btns-more" class="btn-group btn-group-sm" role="group">
                        <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle"
                                data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                            Actions
                        </button>
                        <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                            <a id="suspend" class="dropdown-item ${cls}" href="#">Suspend</a>
                            <a id="resume" class="dropdown-item" href="#">Resume</a>
                            <div class="dropdown-divider"></div>
                            <a id="remove" class="dropdown-item" href="#">Remove</a>
                            <div class="dropdown-divider"></div>
                            <a id="setting" class="dropdown-item" href="#" data-toggle="modal" data-target="#instanceSetModal">Setting</a>
                            <div class="dropdown-divider"></div>
                            <a id="dumpxml" class="dropdown-item" href="${hostUrl}/api/instance/{{uuid}}?format=xml">Dump XML</a>
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
    
        <!-- Virtual Disk -->
        <div id="devices">
        <div id="disk" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseDis" aria-expanded="true" aria-controls="collapseDis">
                    Virtual Disk
                </button>
            </div>
            <div id="collapseDis" class="collapse" aria-labelledby="headingOne" data-parent="#devices">
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
                            <th>Address</th>
                            <th>Format</th>
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
        <!-- Interface -->
        <div id="interface" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseInt" aria-expanded="true" aria-controls="collapseInt">
                    Network Interface
                </button>
            </div>
            <div id="collapseInt" class="collapse" aria-labelledby="headingOne" data-parent="#devices">
            <div class="card-body">
                <div class="card-header-cnt">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#interfaceCreateModal">
                        Attach one
                    </button>
                    <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                    <button id="remove" type="button" class="btn btn-outline-dark btn-sm">Remove</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                </div>
                <div class="">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th><input id="on-all" type="checkbox" aria-label="select all interfaces"></th>
                            <th>ID</th>
                            <th>Model</th>
                            <th>Device</th>
                            <th>Mac</th>
                            <th>Address</th>
                            <th>Source</th>
                        </tr>
                        </thead>
                        <tbody id="display-table">
                        <!-- Loading -->
                        </tbody>
                    </table>
                </div>
            </div>
            </div>
        </div>
        <!-- Graphics -->
        <div id="graphics" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseGra" aria-expanded="true" aria-controls="collapseGra">
                    Graphics Device
                </button>
            </div>
            <div id="collapseGra" class="collapse" aria-labelledby="headingOne" data-parent="#devices">
                <div class="card-body">
                    <div class="card-header-cnt">
                        <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#graphicCreateModal">
                            Attach graphic
                        </button>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                    </div>
                    <div class="">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th><input id="on-all" type="checkbox" aria-label="select all graphics"></th>
                                <th>ID</th>
                                <th>Type</th>
                                <th>Password</th>
                                <th>Listen</th>
                            </tr>
                            </thead>
                            <tbody id="display-table">
                            <!-- Loading -->
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        </div>
        </instance>`)(v);
    }
}
