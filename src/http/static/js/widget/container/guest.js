import {Container} from "./container.js"
import {Utils} from "../../lib/utils.js";
import {GuestCtl} from '../../ctl/guest.js';
import {Api} from "../../api/api.js";
import {InstanceApi} from "../../api/instance.js";

import {DiskCreate} from '../disk/create.js';
import {IsoCreate} from "../disk/iso/create.js";
import {InterfaceCreate} from '../interface/create.js';
import {InstanceSet} from "../instance/setting.js";
import {InstanceRemove} from "../instance/remove.js";
import {GraphicsCreate} from "../graphics/create.js";

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
            this.loading();
        });
    }

    loading() {
        let ctl = new GuestCtl({
            id: this.id(),
            header: {id: this.id("#header")},
            disks: {id: this.id("#disk")},
            interfaces: {id: this.id("#interface")},
            graphics: {id: this.id("#graphics")},
        });
        new InstanceSet({id: this.id('#settingModal'), cpu: ctl.cpu, mem: ctl.mem })
            .onsubmit((e) => {
                ctl.edit(Utils.toJSON(e.form));
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
        let dumpUrl = Api.path(`/api/instance/${v.uuid}?format=xml`);

        return this.compile(`
        <div id="instance" data="{{uuid}}" name="{{name}}" cpu="{{maxCpu}}" memory="{{maxMem}}">
        <div id="header" class="card header">
            <div class="card-header">
                <div class="card-just-left">
                    <a id="refresh" class="none" href="javascript:void(0)">{{name}}</a>
                </div>
            </div>
            <!-- Overview -->
            <div id="collapseOver">
            <div class="card-body">
                <!-- Header buttons -->
                <div class="card-body-hdl">
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
                        </div>
                    </div>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm">{{'refresh' | i}}</button>
                    <div id="power-btns" class="btn-group btn-group-sm" role="group">
                        <button id="start" type="button" class="btn btn-outline-dark btn-sm">
                            {{'power on' | i}}
                        </button>
                        <button id="power" type="button"
                                class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
                                data-toggle="dropdown" aria-expanded="false" ${disabled}>
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
                            <a id="resume" class="dropdown-item ${disabled}" href="javascript:void(0)">{{'resume' | i}}</a>
                            <a id="setting" class="dropdown-item" href="javascript:void(0)" 
                                data-toggle="modal" data-target="#settingModal">{{'setting' | i}}</a>
                            <a id="dumpxml" class="dropdown-item" href="${dumpUrl}">{{'dump xml' | i}}</a>
                            <div class="dropdown-divider"></div>
                            <a id="remove" class="dropdown-item" href="javascript:void(0)" 
                                data-toggle="modal" data-target="#removeModal">{{'remove' | i}}</a>
                        </div>
                    </div>
                </div>
                <div class="card-body-tbl">
                    <div class="overview">
                        <dl class="dl-horizontal">
                            <dt>{{'state' | i}}:</dt>
                            <dd><span class="{{state}}">{{state}}</span></dd>
                            <dt>{{'uuid' | i}}:</dt>
                            <dd>{{uuid}}</dd>
                            <dt>{{'arch' | i}}:</dt>
                            <dd>{{arch}} | {{type}}</dd>
                            <dt>{{'processor' | i}}:</dt>
                            <dd>{{cpuMode ? cpuMode : 'custom'}} | {{maxCpu}} | {{cpuTime}}ms</dd>
                            <dt>{{'memory' | i}}:</dt>
                            <dd>{{maxMem | prettyKiB}} | {{memory | prettyKiB}}</dd>
                        </dl>
                    </div>
                </div>
            </div>
            </div>
        </div>
        
        <div id="collapse">
        <!-- Virtual Disk -->
        <div id="disk" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm" type="button">
                    {{'virtual disk' | i}}
                </button>
            </div>
            <div id="collapseDis">
            <div class="card-body">
                <div class="card-body-hdl">
                    <div id="create-btns" class="btn-group btn-group-sm" role="group">
                        <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#createDiskModal">
                            {{'attach disk' | i}}
                        </button>
                        <button id="creates" type="button"
                            class="btn btn-outline-dark dropdown-toggle dropdown-toggle-split"
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
                    <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
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
        <!-- Interface -->
        <div id="interface" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm" type="button">
                    {{'network interface' | i}}
                </button>
            </div>
            <div id="collapseInt">
            <div class="card-body">
                <div class="card-body-hdl">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#createInterfaceModal">
                        {{'attach interface' | i}}
                    </button>
                    <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                    <button id="remove" type="button" class="btn btn-outline-dark btn-sm">{{'remove' | i}}</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
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
        <!-- Graphics -->
        <!--
        <div id="graphics" class="card device">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm"
                        type="button" data-toggle="collapse"
                        data-target="#collapseGra" aria-expanded="true" aria-controls="collapseGra">
                    {{'graphics device' | i}}
                </button>
            </div>
            <div id="collapseGra" class="collapse" aria-labelledby="headingOne" data-parent="#collapse">
                <div class="card-body">
                    <div class="card-body-hdl">
                        <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                                data-toggle="modal" data-target="#createGraphicModal">
                            {{'attach graphic' | i}}
                        </button>
                        <button id="edit" type="button" class="btn btn-outline-dark btn-sm">{{'edit' | i}}</button>
                        <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    </div>
                    <div class="card-body-tbl">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th><input id="on-all" type="checkbox" aria-label="select all graphics"></th>
                                <th>{{'id' | i}}</th>
                                <th>{{'type' | i}}</th>
                                <th>{{'password' | i}}</th>
                                <th>{{'listen' | i}}</th>
                            </tr>
                            </thead>
                            <tbody id="display-table">
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        -->
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
            <!-- Create ISO/CDROM modal -->
            <div id="createIsoModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create interface modal -->
            <div id="createInterfaceModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
            <!-- Create graphics modal -->
            <div id="createGraphicModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
        </div>
        </div>`, v);
    }
}
