import {Base} from "./base.js"
import {Guest} from "./guest.js"
import {Utils} from "../../com/utils.js";
import {Instances} from '../../ctl/instances.js';
import {Network} from "../../ctl/network.js";
import {Datastores} from "../../ctl/datastores.js";
import {Collapse} from "../collapse.js";
import {Overview} from "../index/overview.js";
import {InstanceCreate} from '../instance/create.js';
import {NATCreate} from "../network/create.js";
import {BridgeCreate} from "../network/bridge/create.js";
import {RoutedCreate} from "../network/routed/create.js";
import {IsolatedCreate} from "../network/isolated/create.js";
import {DirCreate} from "../datastores/create.js";
import {NFSCreate} from "../datastores/nfs/create.js";
import {iSCSICreate} from "../datastores/iscsi/create.js";

export class Index extends Base {
    // {
    //    id: ".container",
    //    default: "instances"
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        super(props);
        this.default = props.default || '/instances';
        console.log('Index', props);
        this.render();
        this.loading();
    }

    loading() {
        this.title('Home - LightStar');
        new Collapse({
            pages: [
                {id: '#collapseSys', name: '/system'},
                {id: '#collapseIns', name: '/instances'},
                {id: '#collapseStore', name: '/datastore'},
                {id: '#collapseNet', name: '/network'}
            ],
            force: this.force,
            default: this.default,
        });
        // loading overview.
        let view = new Overview({
            id: '#overview',
        });
        view.refresh((e) => {
            this.props.name = e.resp.hyper.name;
            $('#system-col').text(this.props.name);
        });
        // register click on overview.
        $('#system-ref').on('click', () => {
            view.refresh();
            $('#collapseSys').collapse('show');
        });

        let ins = new Instances({
            id: '#instances',
            onthis: (e) => {
                console.log("Index.loading", e);
                new Guest({
                    id: this.id,
                    uuid: e.uuid,
                });
            },
        });
        new InstanceCreate({id: '#instanceCreateModal'})
            .onsubmit((e) => {
                ins.create(Utils.toJSON(e.form));
            });
        // loading network.
        let net = new Network({id: '#networks'});
        new NATCreate({id: '#natCreateModal'})
            .onsubmit((e) => {
                net.create(Utils.toJSON(e.form));
            });
        new BridgeCreate({id: '#bridgeCreateModal'})
            .onsubmit((e) => {
                net.create(Utils.toJSON(e.form));
            });
        new RoutedCreate({id: '#routedCreateModal'})
            .onsubmit((e) => {
                net.create(Utils.toJSON(e.form));
            });
        new IsolatedCreate({id: '#isolatedCreateModal'})
            .onsubmit((e) => {
                net.create(Utils.toJSON(e.form));
            });
        // loading datastore.
        let store = new Datastores({
            id: '#datastores',
            upload: '#fileUploadModal',
        });
        new DirCreate({id: '#dirCreateModal'})
            .onsubmit((e) => {
                store.create(Utils.toJSON(e.form));
            });
        new NFSCreate({id: '#nfsCreateModal'})
            .onsubmit((e) => {
                store.create(Utils.toJSON(e.form));
            });
        new iSCSICreate({id: '#iscsiCreateModal'})
            .onsubmit((e) => {
                store.create(Utils.toJSON(e.form));
            });
    }

    template(v) {
        return (`
        <index id="index">
        <!-- System -->
        <div id="system" class="card card-main system">
            <div class="card-header">
                <div class="">
                    <a id="system-col" href="#" data-toggle="collapse"
                       data-target="#collapseSys" aria-expanded="true" aria-controls="collapseSys">${this.props.name}</a>
                    <a class="btn-spot float-right" id="system-ref" href="#system"></a>
                </div>
            </div>
            <div id="collapseSys" class="collapse" aria-labelledby="headingOne" data-parent="#index">
            <div id="overview" class="card-body">
            <!-- Loading -->
            </div>
            </div>
        </div>
        
        <!-- Instances -->
        <div id="instances" class="card instances">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseIns" aria-expanded="true" aria-controls="collapseIns">
                    Guest Instances
                </button>
            </div>
            <div id="collapseIns" class="collapse" aria-labelledby="headingOne" data-parent="#index">
            <div class="card-body">
                <!-- Instances buttons -->
                <div class="card-header-cnt">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#instanceCreateModal">
                        Create new instance
                    </button>
                    <button id="console" type="button" class="btn btn-outline-dark btn-sm">Console</button>
                    <button id="start" type="button" class="btn btn-outline-dark btn-sm">Power on</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                    <button id="more" type="button" class="btn btn-outline-dark btn-sm dropdown-toggle"
                            data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        Actions
                    </button>
                    <div name="btn-more" class="dropdown-menu">
                        <a id="more-start" class="dropdown-item" href="#">Power on</a>
                        <a id="more-shutdown" class="dropdown-item" href="#">Power off</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-reset" class="dropdown-item" href="#">Reset</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-suspend" class="dropdown-item" href="#">Suspend</a>
                        <a id="more-resume" class="dropdown-item" href="#">Resume</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-destroy" class="dropdown-item" href="#">Destroy</a>
                    </div>
                </div>
    
                <!-- Instances display -->
                <div class="">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th>
                                <input id="on-all" type="checkbox" aria-label="select all instances">
                            </th>
                            <th>ID</th>
                            <th>UUID</th>
                            <th>CPU Time</th>
                            <th>Name</th>
                            <th>CPU</th>
                            <th>Memory</th>
                            <th>State</th>
                        </tr>
                        </thead>
                        <tbody id="display-body">
                        <!-- Loading... -->
                        </tbody>
                    </table>
                </div>
            </div>
            </div>
        </div>
        <!-- DataStore -->
        <div id="datastores" class="card card-main">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseStore" aria-expanded="true" aria-controls="collapseStore">
                    Local DataStores
                </button>
            </div>
            <div id="collapseStore" class="collapse" aria-labelledby="headingOne" data-parent="#index">
                <div class="card-body">
                    <!-- DataStore buttons -->
                    <div class="card-header-cnt">
                        <div id="create-btns" class="btn-group btn-group-sm" role="group">
                            <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                    data-toggle="modal" data-target="#dirCreateModal">
                                New a datastore
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-nfs" class="dropdown-item" data-toggle="modal" data-target="#nfsCreateModal">
                                    New NFS datastore
                                </a>
                                <div class="dropdown-divider"></div>
                                <a id="create-iscsi" class="dropdown-item" data-toggle="modal" data-target="#iscsiCreateModal">
                                    New iSCSI datastore
                                </a>
                            </div>
                        </div>
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#fileUploadModal">
                            Upload
                        </button>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                        <button id="delete" type="button" class="btn btn-outline-dark btn-sm">Delete</button>
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                    </div>
    
                    <!-- DataStore display -->
                    <div class="l-content-body">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th><input id="on-all" type="checkbox"></th>
                                <th>ID</th>
                                <th>UUID</th>
                                <th>Name</th>
                                <th>Source</th>
                                <th>Capacity</th>
                                <th>Available</th>
                                <th>State</th>
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
    
        <!-- Network -->
        <div id="networks" class="card card-main">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseNet" aria-expanded="true" aria-controls="collapseNet">
                    Virtual Networks
                </button>
            </div>
            <div id="collapseNet" class="collapse" aria-labelledby="headingOne" data-parent="#index">
                <div class="card-body">
                    <!-- Network buttons -->
                    <div class="card-header-cnt">
                        <div id="create-btns" class="btn-group btn-group-sm" role="group">
                            <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                    data-toggle="modal" data-target="#natCreateModal">
                                Create network
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-routed" class="dropdown-item" data-toggle="modal" data-target="#routedCreateModal">
                                    Routed network
                                </a>
                                <div class="dropdown-divider"></div>
                                <a id="create-isolated" class="dropdown-item" data-toggle="modal" data-target="#isolatedCreateModal">
                                    Isolated network
                                </a>
                                <div class="dropdown-divider"></div>
                                <a id="create-bridge" class="dropdown-item" data-toggle="modal" data-target="#bridgeCreateModal">
                                    Host bridge
                                </a>
                            </div>
                        </div>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                        <button id="delete" type="button" class="btn btn-outline-dark btn-sm">Delete</button>
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >Refresh</button>
                    </div>
    
                    <!-- Network display -->
                    <div class="l-content-body">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th><input id="on-all" type="checkbox"></th>
                                <th>ID</th>
                                <th>UUID</th>
                                <th>Name</th>
                                <th>Address</th>
                                <th>Mode</th>
                                <th>State</th>
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
        </index>`)
    }
}
