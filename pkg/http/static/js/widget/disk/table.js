import {Widget} from "../widget.js";
import {DiskApi} from "../../api/disk.js";


export class DiskTable extends Widget {
    // {
    //   id: '#xx',
    //   inst: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
        this.name = props.name;
    }

    loading() {
        return `<tr><td colspan="8" class="text-center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new DiskApi({
            tasks: this.tasks,
            inst: this.inst,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        {{if (items.length === 0)}}
            <tr>
                <td colspan="8" class="text-center">{{'no data to display' | i}}</td>
            </tr>
        {{/if}}
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.device}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.bus}}</td>
                <td>{{v.device}}</td>
                <td>{{v.source}}</td>
                <td>{{if v.volume.type === ""}} - {{else}} {{v.volume.capacity | prettyByte}} {{/if}}</td>
                <td>{{if v.volume.type === ""}} - {{else}} {{v.volume.allocation | prettyByte}} {{/if}}</td>
                <td><span>
                    {{if v.addrType == "pci"}}
                        pci:{{v.addrBus}}:{{v.addrSlot}}.{{v.addrFunc}}
                    {{else if  v.addrType == "drive"}}
                        drv:{{v.addrBus}}:{{v.addrTgt}}.{{v.addrUnit}}
                    {{/if}}</span>
                </td>
            </tr>
        {{/each}}
        `, data);
    }
}
