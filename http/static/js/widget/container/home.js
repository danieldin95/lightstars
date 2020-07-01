import {Container} from "./container.js"
import {Guest} from "./guest.js"
import {Network} from "./network.js";
import {Utils} from "../../com/utils.js";
import {Location} from "../../com/location.js";
import {InstanceCtl} from '../../ctl/instance.js';
import {NetworksCtl} from "../../ctl/networks.js";
import {DataStoresCtl} from "../../ctl/datastores.js";
import {Collapse} from "../collapse.js";
import {Overview} from "../index/overview.js";
import {InstanceCreate} from '../instance/create.js';
import {NATCreate} from "../network/create.js";
import {BridgeCreate} from "../network/bridge/create.js";
import {IsolatedCreate} from "../network/isolated/create.js";
import {DirCreate} from "../datastore/create.js";
import {NFSCreate} from "../datastore/nfs/create.js";
import {iSCSICreate} from "../datastore/iscsi/create.js";
import {Pool} from "./pool.js";

export class Home extends Container {
    // {
    //    parent: "#Container",
    //    default: "/instances"
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        super(props);
        this.current = '#index';
        this.default = props.default || '/instances';
        console.log('Index', props);
        this.render();
        this.loading();
    }

    loading() {
        this.title('Home');
        new Collapse({
            pages: [
                {id: this.id('#collapseSys'), name: '/system'},
                {id: this.id('#collapseIns'), name: '/instances'},
                {id: this.id('#collapseStore'), name: '/datastore'},
                {id: this.id('#collapseNet'), name: '/network'}
            ],
            force: this.force,
            default: this.default,
        });
        // loading overview.
        let view = new Overview({
            id: this.id('#overview'),
        });
        view.refresh((e) => {
            this.props.name = e.resp.hyper.name;
            $(this.id('#system-col')).text(this.props.name);
        });
        // register click on overview.
        $(this.id('#system-ref')).on('click', () => {
            view.refresh();
            $(this.id('#collapseSys')).collapse('show');
        });

        let iCtl = new InstanceCtl({
            id: this.id('#instances'),
            onthis: (e) => {
                console.log("Guest.loading", e);
                new Guest({
                    parent: this.parent,
                    uuid: e.uuid,
                });
            },
        });
        new InstanceCreate({id: '#createGuestModal'})
            .onsubmit((e) => {
                iCtl.create(Utils.toJSON(e.form));
            });
        // loading network.
        let nCtl = new NetworksCtl({
            id: this.id('#networks'),
            onthis: (e) => {
                console.log("network.loading", e);
                new Network({
                    parent: this.parent,
                    uuid: e.uuid,
                });
            },
        });
        new NATCreate({id: '#createNatModal'})
            .onsubmit((e) => {
                nCtl.create(Utils.toJSON(e.form));
            });
        new BridgeCreate({id: '#createBridgeModal'})
            .onsubmit((e) => {
                nCtl.create(Utils.toJSON(e.form));
            });
        new IsolatedCreate({id: '#createIsolatedModal'})
            .onsubmit((e) => {
                nCtl.create(Utils.toJSON(e.form));
            });
        // loading data storage.
        let sCtl = new DataStoresCtl({
            id: this.id('#datastores'),
            onthis: (e) => {
                console.log("storage loading", e)
                new Pool({
                    parent: this.parent,
                    uuid: e.uuid,
                });
            },
            upload: '#uploadFileModal',
        });
        new DirCreate({id: '#createDirModal'})
            .onsubmit((e) => {
                sCtl.create(Utils.toJSON(e.form));
            });
        new NFSCreate({id: '#createNfsModal'})
            .onsubmit((e) => {
                sCtl.create(Utils.toJSON(e.form));
            });
        new iSCSICreate({id: '#createIscsiModal'})
            .onsubmit((e) => {
                sCtl.create(Utils.toJSON(e.form));
            });
    }

    template(v) {
        let query = Location.query();
        return (`
        <div id="index">
        <!-- System -->
        <div id="system" class="card card-main system">
            <div class="card-header">
                <div class="">
                    <a id="system-col" href="javascript:void(0)" data-toggle="collapse"
                       data-target="#collapseSys" aria-expanded="true" aria-controls="collapseSys">${this.props.name}</a>
                    <a class="btn-spot float-right" id="system-ref" href="#/system?${query}"></a>
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
                            data-toggle="modal" data-target="#createGuestModal">
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
                        <a id="more-start" class="dropdown-item" href="javascript:void(0)">Power on</a>
                        <a id="more-shutdown" class="dropdown-item" href="javascript:void(0)">Power off</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-suspend" class="dropdown-item" href="javascript:void(0)">Suspend</a>
                        <a id="more-resume" class="dropdown-item" href="javascript:void(0)">Resume</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-reset" class="dropdown-item" href="javascript:void(0)">Reset</a>                        
                        <a id="more-destroy" class="dropdown-item" href="javascript:void(0)">Destroy</a>
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
                                    data-toggle="modal" data-target="#createDirModal">
                                New a datastore
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-nfs" class="dropdown-item" data-toggle="modal" data-target="#createNfsModal">
                                    New nfs datastore
                                </a>
                                <a id="create-iscsi" class="dropdown-item" data-toggle="modal" data-target="#createIscsiModal">
                                    New iscsi datastore
                                </a>
                            </div>
                        </div>
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#uploadFileModal">
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
                                    data-toggle="modal" data-target="#createNatModal">
                                Create network
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-isolated" class="dropdown-item" data-toggle="modal" data-target="#createIsolatedModal">
                                    Create isolated network
                                </a>
                                <a id="create-bridge" class="dropdown-item" data-toggle="modal" data-target="#createBridgeModal">
                                    Create existing bridge
                                </a>
                            </div>
                        </div>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">Edit</button>
                        <button id="delete" type="button" class="btn btn-outline-dark btn-sm">Remove</button>
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
                                <th>Bridge</th>
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
        
        <!-- Modal -->
        <div id="modals">
            <!-- Create instance modal -->
            <div id="createGuestModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create datastore modal -->
            <div id="createDirModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Upload file modal -->
            <div id="uploadFileModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createNfsModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createIscsiModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create network modal -->
            <div id="createNatModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createBridgeModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createIsolatedModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>        
        </div>
        </div>`)
    }
}
