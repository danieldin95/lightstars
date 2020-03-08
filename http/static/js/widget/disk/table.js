import {DiskApi} from "../../api/disk.js";


export class DiskTable {
    // {
    //   id: '#xx',
    //   instance: 'uuid',
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
        this.instance = props.instance;
    }

    loading() {
        return `<tr><td colspan="7" style="text-align: center">Loading...</td></tr>`
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new DiskApi({
            tasks: this.tasks,
            instance: this.instance,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return template.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.device}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.bus}}</td>
                <td>{{v.device}}</td>
                <td>{{v.source}}</td>
                <td><span>
                    {{if v.addrType == "pci"}}
                        pci:{{v.addrBus | aton 2}}:{{v.addrSlot | aton 2}}.{{v.addrFunc}}
                    {{else if  v.addrType == "drive"}}
                        drv:{{v.addrBus | aton 2}}:{{v.addrTgt | aton 2}}.{{v.addrUnit}}
                    {{/if}}</span>
                </td>
                <td>{{v.format}}</td>
            </tr>
        {{/each}}
        `)(data)
    }
}