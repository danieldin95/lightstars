import {Container} from "./container.js"
import {Network} from "./network.js";
import {Utils} from "../lib/utils.js";
import {NetworksCtl} from "../controller/networks.js";
import {NATCreate} from "../widget/network/create.js";
import {OVSCreate} from "../widget/network/ovs/create.js";
import {BridgeCreate} from "../widget/network/bridge/create.js";
import {IsolatedCreate} from "../widget/network/isolated/create.js";
import {I18N} from "../lib/i18n.js";

export class Networks extends Container {
    // {
    //    parent: "#container",
    // }
    constructor(props) {
        super(props);
        this.current = '#index';

        this.render();
        this.loading();
    }

    loading() {
        this.title(I18N.i('network'));
        // loading network.
        let nCtl = new NetworksCtl({
            id: this.id('#networks'),
            onthis: (e) => {
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
        new OVSCreate({id: '#createOvsModal'})
            .onsubmit((e) => {
                nCtl.create(Utils.toJSON(e.form));
            });
    }

    template(v) {
        return this.compile(`
        <div id="index">
        <!-- Network -->
        <div id="networks" class="card shadow">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm" type="button">
                    {{'virtual networks' | i}}
                </button>
            </div>
            <div class="card-body">
                <!-- Network buttons -->
                <div class="row card-body-hdl">
                    <div class="col-auto mr-auto">
                        <div id="create-btns" class="btn-group btn-group-sm" role="group">
                            <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                    data-toggle="modal" data-target="#createNatModal">
                                {{'create network' | i}}
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-info dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-bridge" class="dropdown-item" data-toggle="modal" data-target="#createBridgeModal">
                                    {{'linux bridge based' | i}}
                                </a>                            
                                <a id="create-ovs" class="dropdown-item" data-toggle="modal" data-target="#createOvsModal">
                                    {{'open vswitch based' | i}}
                                </a>
                            </div>
                        </div>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                        <button id="delete" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
                    </div>
                    <div class="col-auto">
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    </div>
                </div>

                <!-- Network display -->
                <div class="card-body-tbl">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th><input id="on-all" type="checkbox"></th>
                            <th>{{'id' | i}}</th>
                            <th>{{'name' | i}}</th>
                            <th>{{'bridge' | i}}</th>
                            <th>{{'state' | i}}</th>
                        </tr>
                        </thead>
                        <tbody id="display-table">
                        <!-- Loading -->
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <!-- Modal -->
        <div id="modals">
            <!-- Create network modal -->
            <div id="createNatModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createBridgeModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createIsolatedModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createOvsModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>          
        </div>
        </div>`)
    }
}
