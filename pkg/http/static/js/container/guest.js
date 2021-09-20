import {Container} from "./container.js"
import {Utils} from "../lib/utils.js";
import {GuestCtl} from '../controller/guest.js';
import {Api} from "../api/api.js";
import {InstanceApi} from "../api/instance.js";

import {DiskCreate} from '../widget/disk/create.js';
import {IsoCreate} from "../widget/disk/iso/create.js";
import {InterfaceCreate} from '../widget/interface/create.js';
import {InstanceSet} from "../widget/instance/setting.js";
import {InstanceRemove} from "../widget/instance/remove.js";
import {GraphicsCreate} from "../widget/graphics/create.js";
import {TitleSet} from "../widget/instance/title.js";
import {SnapshotCreate} from "../widget/snapshot/create.js";
import {DiskRemove} from "../widget/disk/remove.js";
import {DiskApi} from "../api/disk.js";

export class Guest extends Container {
    // {
    //    parent: "#container",
    //    uuid: "",
    //    default: "disk"
    // }
    constructor(props) {
        super(props);
        this.current = "#instance";
        this.name = "";
        this.uuid = props.uuid;

        this.render();
    }

    render() {
        new InstanceApi({uuids: this.uuid}).get(this, (e) => {
            this.name = e.resp.name;
            this.title(this.name);
            this.view = $(this.template(e.resp));
            this.view.find('#header #refresh').on('click', (e) => {
                this.render();
            });
            $(this.parent).html(this.view);
            this.loading(e.resp);
        });
    }

    loading(data) {
        let ctl = new GuestCtl({
            id: this.id(),
            name: data.name,
            uuid: data.uuid,
            header: {id: this.id("#header")},
            disks: {
                id: this.id("#disk"),
                onRemove: (objs) => {
                    new DiskRemove({
                        id: this.id('#removeDiskModal'),
                        name: objs.uuids,
                    }).onsubmit((e) => {
                        new DiskApi(objs).delete();
                    });
                }
            },
            interfaces: {id: this.id("#interface")},
            graphics: {id: this.id("#graphics")},
            snapshot: {id: this.id('#snapshot')},
            data: data,
        });
        new InstanceSet({id: this.id('#settingModal'), data: data })
            .onsubmit((e) => {
                ctl.edit(Utils.toJSON(e.form));
            });
        new TitleSet({id: this.id('#settingTitleModal'), data: data })
            .onsubmit((e) => {
                ctl.title(Utils.toJSON(e.form));
            });
        new InstanceRemove({id: this.id('#removeModal'), name: this.name, uuid: this.uuid })
            .onsubmit((e) => {
                ctl.remove();
            });
        // loading disks and interfaces.
        new DiskCreate({id: this.id('#createDiskModal')})
            .onsubmit((e) => {
                ctl.disk.create(Utils.toJSON(e.form));
            });
        new IsoCreate({id: this.id("#createIsoModal")})
            .onsubmit((e) => {
                ctl.disk.create(Utils.toJSON(e.form));
            });
        new InterfaceCreate({id: this.id('#createInterfaceModal')})
            .onsubmit((e) => {
                ctl.interface.create(Utils.toJSON(e.form));
            });
        new GraphicsCreate({id: this.id('#createGraphicModal')})
            .onsubmit((e) => {
                ctl.graphics.create(Utils.toJSON(e.form));
            });
        new SnapshotCreate({id: this.id('#createSnapshotModal')})
            .onsubmit((e) => {
                ctl.snapshot.create(Utils.toJSON(e.form));
            });
        // register console draggable.
        $((e) => {
            $(this.id('#consoleModal')).draggable();
        });
    }

    template(v) {
        let disabled = v.state === 'running' ? '' : 'disabled';
        let host = Api.host();
        let pass = Utils.graphic(v, 'vnc', 'password');
        let vncUrl = `/ui/console?id=${v.uuid}&password=${pass}&node=${host}&title=${v.name}`;
        let liteUrl = `/ui/lite?id=${v.uuid}&password=${pass}&node=${host}&title=${v.name}`;
        let dumpUrl = Api.path(`/api/instance/${v.uuid}?format=xml`);
        let os = Utils.os();
        let localUrl = Api.path(`/api/instance/${v.uuid}/graphics?format=vv&os=${os}`);

        return this.compile(`
        <div id="instance" data="{{uuid}}" name="{{name}}">
        <div id="header" class="card header shadow">
            <div class="card-header">
                <div class="text-left">
                    <a id="refresh" class="none" href="javascript:void(0)">{{name}}</a>
                </div>
            </div>
            <!-- Overview -->
            <div id="collapseOver">
            <div class="card-body">
                <!-- Header buttons -->
                <div class="row card-body-hdl">
                    <div class="col-auto mr-auto">
                        <div id="console-btns" class="btn-group btn-group-sm" role="group">
                            <button id="console" type="button" class="btn btn-outline-dark btn-sm"
                                    data-toggle="modal" data-target="#consoleModal" data="${vncUrl}" ${disabled}>{{'console' | i}}</button>
                            <button id="consoles" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false" ${disabled}>
                                <span class="sr-only">Toggle Dropdown</span></button>
                            <div id="console-more" class="dropdown-menu" aria-labelledby="consoles">
                                <a id="console-self" class="dropdown-item" href="javascript:void(0)" data="${vncUrl}">
                                    {{'console in self' | i}}
                                </a>
                                <a id="console-blank" class="dropdown-item" href="javascript:void(0)" data="${vncUrl}">
                                    {{'console in new blank' | i}}
                                </a>
                                <a id="console-window" class="dropdown-item" href="javascript:void(0)" data="${vncUrl}">
                                    {{'console in new window' | i}}
                                </a>
                                <a id="console-local" class="dropdown-item" href="${localUrl}">
                                    {{'console by remote viewer' | i}}
                                </a>
                            </div>
                        </div>
                        <div id="power-btns" class="btn-group btn-group-sm" role="group">
                            <button id="start" type="button" class="btn btn-outline-dark btn-sm">
                                {{'power on' | i}}
                            </button>
                            <button id="power" type="button"
                                    class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                    data-toggle="dropdown" aria-expanded="false">
                                <span class="sr-only">Toggle Dropdown</span>
                            </button>
                            <div id="power-more" class="dropdown-menu" aria-labelledby="power">
                                <a id="start" class="dropdown-item" href="javascript:void(0)">{{'power on' | i}}</a>
                                <a id="shutdown" class="dropdown-item" href="javascript:void(0)">{{'power off' | i}}</a>
                                <div class="dropdown-divider"></div>
                                <a id="reset" class="dropdown-item" href="javascript:void(0)">{{'reset' | i}}</a>
                                <a id="destroy" class="dropdown-item" href="javascript:void(0)">{{'destroy' | i}}</a>
                            </div>
                        </div>
                        <div id="btns-more" class="btn-group btn-group-sm" role="group">
                            <button id="btns-more" type="button" class="btn btn-outline-dark dropdown-toggle btn-sm"
                                    data-toggle="dropdown" aria-expanded="true" aria-expanded="false">
                                {{'actions' | i}}
                            </button>
                            <div name="btn-more" class="dropdown-menu" aria-labelledby="btns-more">
                                <a id="suspend" class="dropdown-item ${disabled}" href="javascript:void(0)">{{'suspend' | i}}</a>
                                <a id="resume" class="dropdown-item" href="javascript:void(0)">{{'resume' | i}}</a>
                                <a id="setting" class="dropdown-item" href="javascript:void(0)" 
                                    data-toggle="modal" data-target="#settingModal">{{'setting spec' | i}}</a>
                                <a id="title" class="dropdown-item" href="javascript:void(0)" 
                                    data-toggle="modal" data-target="#settingTitleModal">{{'setting title' | i}}</a>                                    
                                <a id="dumpxml" class="dropdown-item" href="${dumpUrl}">{{'dump xml' | i}}</a>
                                <div class="dropdown-divider"></div>
                                <a id="remove" class="dropdown-item" href="javascript:void(0)" 
                                    data-toggle="modal" data-target="#removeModal">{{'remove' | i}}</a>
                            </div>
                        </div>
                    </div>
                    <div class="col-auto">
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">{{'refresh' | i}}</button>
                    </div>
                </div>
                <div class="card-body-tbl row">
                    <div class="col-sm-12 col-md-5 col-lg-4 mt-1 pt-3 split-vertical">
                        <div style="width: 328px; height: 188px; background-color: rgb(40 40 40); border-radius: 4px; padding: 4px;">
                            <iframe width="320px" height="180px" src="${liteUrl}" frameborder="0"></iframe>
                        </div>
                    </div>
                    <div class="col-sm-12 col-md-7 col-lg-8 mt-1 split-vertical">
                        <dl class="dl-horizontal dl-horizontal-r">
                            <dt>{{'name' | i}}:</dt>
                            <dd>&nbsp;{{name}}</dd>
                            <dt>{{'state' | i}}:</dt>
                            <dd>&nbsp;<span class="st-{{state}}">{{state}}</span></dd>
                            <dt>{{'uuid' | i}}:</dt>
                            <dd>&nbsp;{{uuid}}</dd>
                            <dt>{{'title' | i}}:</dt>
                            <dd>&nbsp;{{title}}</dd>                            
                            <dt>{{'arch' | i}}:</dt>
                            <dd>&nbsp;{{arch}} | {{type}}</dd>
                            <dt>{{'processor' | i}}:</dt>
                            <dd title="{{'model | number | time' | i}}">
                              &nbsp;{{cpuMode | prettyCpuMode}} | {{maxCpu}} | {{cpuTime}}ms
                            </dd>
                            <dt>{{'memory' | i}}:</dt>
                            <dd title="{{'max | current' | i}}">
                              &nbsp;{{maxMem | prettyKiB}} | {{memory | prettyKiB}}
                            </dd>
                        </dl>
                    </div>
                </div>
            </div>
            </div>
        </div>
        
        <div class="card-tab">
            <ul class="nav nav-pills justify-content-start" id="pills-tab" role="tablist">
              <li class="nav-item" role="presentation">
                <a class="nav-link active" id="pills-0-tab" data-toggle="pill" href="#pills-0" 
                    role="tab" aria-controls="pills-0" aria-selected="true">{{'virtual disk' | i}}</a>
              </li>
              <li class="nav-item" role="presentation">
                <a class="nav-link" id="pills-1-tab" data-toggle="pill" href="#pills-1" 
                    role="tab" aria-controls="pills-1" aria-selected="false">{{'network interface' | i}}</a>
              </li>
              <li class="nav-item" role="presentation">
                <a class="nav-link" id="pills-2-tab" data-toggle="pill" href="#pills-2" 
                    role="tab" aria-controls="pills-2" aria-selected="false">{{'instance snapshot' | i}}</a>
              </li>
            </ul>
            <div class="tab-content" id="pills-tabContent">
              <div class="tab-pane fade show active" id="pills-0" role="tabpanel" aria-labelledby="pills-0-tab">
                <!-- Virtual Disk -->
                <div id="disk" class="card shadow">
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <div id="create-btns" class="btn-group btn-group-sm" role="group">
                                    <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                            data-toggle="modal" data-target="#createDiskModal">
                                        {{'attach disk' | i}}
                                    </button>
                                    <button id="creates" type="button"
                                        class="btn btn-outline-info dropdown-toggle dropdown-toggle-split"
                                        data-toggle="dropdown" aria-expanded="false">
                                        <span class="sr-only">Toggle Dropdown</span>
                                    </button>
                                    <div id="create-more" class="dropdown-menu" aria-labelledby="creates">
                                        <a id="create-iso" class="dropdown-item" data-toggle="modal" data-target="#createIsoModal">
                                            {{'attach cdrom' | i}}
                                        </a>
                                    </div>
                                </div>  
                                <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                                <button id="remove" type="button" class="btn btn-outline-dark btn-sm"
                                    data-toggle="modal" data-target="#removeDiskModal">{{'remove' | i}}</button>
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
                                    <th>{{'bus' | i}}</th>
                                    <th>{{'device' | i}}</th>
                                    <th>{{'source' | i}}</th>
                                    <th>{{'capacity' | i}}</th>
                                    <th>{{'allocation' | i}}</th>
                                    <th>{{'address' | i}}</th>
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
                <!-- Interface -->
                <div id="interface" class="card shadow">
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                        data-toggle="modal" data-target="#createInterfaceModal">
                                    {{'attach interface' | i}}
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
                                    <th><input id="on-all" type="checkbox" aria-label="select all interfaces"></th>
                                    <th>{{'id' | i}}</th>
                                    <th>{{'model' | i}}</th>
                                    <th>{{'device' | i}}</th>
                                    <th>{{'mac' | i}}</th>
                                    <th>{{'address' | i}}</th>
                                    <th>{{'source' | i}}</th>
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
              <div class="tab-pane fade" id="pills-2" role="tabpanel" aria-labelledby="pills-2-tab">
                <!-- Snapshots -->
                <div id="snapshot" class="card shadow">
                    <div class="card-body">
                        <div class="row card-body-hdl">
                            <div class="col-auto mr-auto">
                                <button id="create" type="button" class="btn btn-outline-info btn-sm"
                                        data-toggle="modal" data-target="#createSnapshotModal">
                                    {{'create snapshot' | i}}
                                </button>
                                <button id="revert" type="button" class="btn btn-outline-dark btn-sm">{{'revert' | i}}</button>
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
                                    <th><input id="on-all" type="checkbox" aria-label="select all"></th>
                                    <th>{{'id' | i}}</th>
                                    <th>{{'name' | i}}</th>
                                    <th>{{'uptime' | i}}</th>
                                    <th>{{'state' | i}}</th>
                                </tr>
                                </thead>
                                <tbody id="display-table">
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
              </div>
            </div>
        </div>
        
        <!-- Modals -->
        <div id="modals">
            <!-- Console modal -->
            <div id="consoleModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true">
                <div class="modal-dialog modal-lg modal-console" role="document">
                    <div class="modal-content">
                        <div class="modal-body"></div>
                    </div>
                </div>
            </div>
            <!-- Remove confirm -->
            <div id="removeModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Setting instance modal -->
            <div id="settingModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create disk modal -->
            <div id="createDiskModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <div id="removeDiskModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create ISO/CDROM modal -->
            <div id="createIsoModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create interface modal -->
            <div id="createInterfaceModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create snapshot modal -->
            <div id="createSnapshotModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Setting title modal -->
            <div id="settingTitleModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        </div>`, v);
    }
}
