import {Container} from "./container.js"
import {Guest} from "./guest.js"
import {Utils} from "../../lib/utils.js";
import {InstanceCtl} from '../../ctl/instance.js';
import {InstanceCreate} from '../instance/create.js';
import {I18N} from "../../lib/i18n.js";
import {InstanceApi} from "../../api/instance.js";

export class Instances extends Container {
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
        this.title(I18N.i('instances'));

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
                new InstanceApi().create(Utils.toJSON(e.form));
            });
    }

    template(v) {
        return this.compile(`
        <div id="index">
        
        <!-- Instances -->
        <div id="instances" class="card instances">
            <div class="card-header">
                <button class="btn btn-link btn-block text-left btn-sm" type="button">
                    {{'guest instances' | i}}
                </button>
            </div>
            <div id="collapseIns">
            <div class="card-body">
                <!-- Instances buttons -->
                <div class="card-body-hdl">
                    <button id="create" type="button" class="btn btn-outline-dark btn-sm"
                            data-toggle="modal" data-target="#createGuestModal">
                        {{'create new instance' | i}}
                    </button>
                    <button id="console" type="button" class="btn btn-outline-dark btn-sm">{{'console' | i}}</button>
                    <button id="start" type="button" class="btn btn-outline-dark btn-sm">{{'power on' | i}}</button>
                    <button id="refresh" type="button" class="btn btn-outline-dark btn-sm" >{{'refresh' | i}}</button>
                    <button id="more" type="button" class="btn btn-outline-dark btn-sm dropdown-toggle"
                            data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        {{'actions' | i}}
                    </button>
                    <div name="btn-more" class="dropdown-menu">
                        <a id="more-start" class="dropdown-item" href="javascript:void(0)">{{'power on' | i}}</a>
                        <a id="more-shutdown" class="dropdown-item" href="javascript:void(0)">{{'power off' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-suspend" class="dropdown-item" href="javascript:void(0)">{{'suspend' | i}}</a>
                        <a id="more-resume" class="dropdown-item" href="javascript:void(0)">{{'resume' | i}}</a>
                        <div class="dropdown-divider"></div>
                        <a id="more-reset" class="dropdown-item" href="javascript:void(0)">{{'reset' | i}}</a>                        
                        <a id="more-destroy" class="dropdown-item" href="javascript:void(0)">{{'destroy' | i}}</a>
                    </div>
                </div>
    
                <!-- Instances display -->
                <div class="card-body-tbl">
                    <table class="table table-striped">
                        <thead>
                        <tr>
                            <th>
                                <input id="on-all" type="checkbox" aria-label="select all instances">
                            </th>
                            <th>{{'id' | i}}</th>
                            <th>{{'uuid' | i}}</th>
                            <th>{{'name' | i}}</th>
                            <th>{{'processor' | i}}</th>
                            <th>{{'memory' | i}}</th>
                            <th>{{'state' | i}}</th>
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
        
        <!-- Modal -->
        <div id="modals">
            <!-- Create instance modal -->
            <div id="createGuestModal" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>      
        </div>
        </div>`)
    }
}
