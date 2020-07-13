import {Widget} from "../widget.js";
import {InterfaceApi} from "../../api/interface.js";


export class InterfaceTable extends Widget {
    // {
    //   id: '#xx',
    //   inst: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
    }

    loading() {
        return `<tr><td colspan="7" style="text-align: center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new InterfaceApi({
            tasks: this.tasks,
            inst: this.inst,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.address}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.model}}</td>
                <td>{{v.device}}</td>
                <td>{{v.address}}</td>
                <td><span>
                {{if  v.addrType == "pci"}}
                    pci:{{v.addrBus | a2n 2}}:{{v.addrSlot | a2n 2}}.{{v.addrFunc}}
                {{/if}}</span>
                </td>
                <td>{{v.source == "" ? v.network : v.source}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
