import {Container} from "./container.js"
import {Utils} from "../../lib/utils.js";
import {DataStoresCtl} from "../../ctl/datastores.js";
import {DirCreate} from "../datastore/create.js";
import {NFSCreate} from "../datastore/nfs/create.js";
import {iSCSICreate} from "../datastore/iscsi/create.js";
import {Pool} from "./pool.js";
import {I18N} from "../../lib/i18n.js";

export class DataStores extends Container {
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
        this.title(I18N.i('datastore'));
        // loading data storage.
        let sCtl = new DataStoresCtl({
            id: this.id('#datastores'),
            onthis: (e) => {
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
        return this.compile(`
        <div id="index">
        <!-- DataStore -->
        <div id="datastores" class="card card-main">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm" type="button">
                    {{'local datastores' | i}}
                </button>
            </div>
            <div id="collapseStore">
                <div class="card-body">
                    <!-- DataStore buttons -->
                    <div class="card-body-hdl">
                        <div id="create-btns" class="btn-group btn-group-sm" role="group">
                            <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                    data-toggle="modal" data-target="#createDirModal">
                                {{'new a datastore' | i}}
                            </button>
                            <button id="creates" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                <a id="create-nfs" class="dropdown-item" data-toggle="modal" data-target="#createNfsModal">
                                    {{'new nfs datastore' | i}}
                                </a>
                                <a id="create-iscsi" class="dropdown-item" data-toggle="modal" data-target="#createIscsiModal">
                                    {{'new iscsi datastore' | i}}
                                </a>
                            </div>
                        </div>
                        <button id="upload" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#uploadFileModal">
                            {{'upload' | i}}
                        </button>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                        <button id="delete" type="button" class="btn btn-outline-dark btn-sm">{{'delete' | i}}</button>
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    </div>
    
                    <!-- DataStore display -->
                    <div class="card-body-tbl">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th><input id="on-all" type="checkbox"></th>
                                <th>{{'id' | i}}</th>
                                <th>{{'uuid' | i}}</th>
                                <th>{{'name' | i}}</th>
                                <th>{{'source' | i}}</th>
                                <th>{{'capacity' | i}}</th>
                                <th>{{'allocation' | i}}</th>
                                <th>{{'state' | i}}</th>
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
    
        
        <!-- Modal -->
        <div id="modals">
            <!-- Create datastore modal -->
            <div id="createDirModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Upload file modal -->
            <div id="uploadFileModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createNfsModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="createIscsiModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>    
        </div>
        </div>`)
    }
}
