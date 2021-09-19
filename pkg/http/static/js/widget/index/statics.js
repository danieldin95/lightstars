import {Widget} from "../widget.js";
import {HyperApi} from "../../api/hyper.js";


export class Statics extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
        console.log(props);
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        new HyperApi({tasks: this.tasks}).statics(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            if (func) {
                func({data, resp: e.resp});
            }
        });
    }

    render(data) {
        return this.compile(`
            <dl class="dl-horizontal dl-horizontal-r">
                <dt>{{'datastore' | i}}:</dt>
                <dd>
                    <span class="badge badge-pill badge-lite" title="total">{{datastore.total}}</span>
                    <span class="badge badge-pill badge-lite st-active" title="active">{{datastore.active}}</span>
                </dd>
                <dt>{{'guest instances' | i}}:</dt>
                <dd>
                    <span class="badge badge-pill badge-lite" title="total">{{instance.total}}</span>
                    <span class="badge badge-pill badge-lite st-running" title="running">{{instance.active}}</span>
                    <span class="badge badge-pill badge-lite st-inactive" title="shutdown">{{instance.inactive}}</span>
                </dd>
                <dt>{{'network' | i}}:</dt>
                <dd>
                    <span class="badge badge-pill badge-lite" title="total">{{network.total}}</span>
                    <span class="badge badge-pill badge-lite st-active" title="active">{{network.active}}</span>
                </dd>
                <dt>{{'virtual ports' | i}}:</dt>
                <dd>
                    <span class="badge badge-pill badge-lite" title="total">{{ports.total}}</span>
                    <span class="badge badge-pill badge-lite st-active" title="up">{{ports.active}}</span>
                    <span class="badge badge-pill badge-lite st-inactive" title="down">{{ports.inactive}}</span>
                </dd>
            </dl>`, data);
    }
}
