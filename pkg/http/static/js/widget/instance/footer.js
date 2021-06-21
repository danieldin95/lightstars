import {Widget} from "../widget.js";
import {InstanceApi} from "../../api/instance.js";


export class InstanceFooter extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
        this.api = new InstanceApi({tasks: this.tasks});
    }

    loading() {
        return `<div style="text-align:center">Loading...</div>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        this.api.stats(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        <div class="row">
            <div class="col-auto mr-auto ml-auto">
                <span class="badge badge-pill badge-outline" title="running | shutoff | others">
                    Total {{running}} | {{shutdown}} | {{others}}
                </span>
                <span class="badge badge-pill badge-outline" title="occupied | alloc">
                    Memory {{occupiedMem | prettyKiB}} | {{allocMem | prettyKiB}}
                </span>
                <span class="badge badge-pill badge-outline" title="occupied | alloc">
                    CPU {{occupiedCpu}} | {{allocCpu}}
                </span>
                <span class="badge badge-pill badge-outline" title="size">
                    Storage {{allocStorage | prettyByte}}
                </span>
            </div>
        </div>`, data);
    }
}
