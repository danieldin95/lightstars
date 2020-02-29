import {InterfaceApi} from "../../api/interface.js";


export class InterfaceTable {
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
        $(this.id).html(this.loading());
        console.log("InterfaceTable.refresh", data, func);
        new InterfaceApi({
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
                <td><input id="on-one" type="checkbox" data="{{v.address}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.model}}</td>
                <td>{{v.device}}</td>
                <td>{{v.address}}</td>
                <td><span>
                {{if  v.addrType == "pci"}}
                    pci:{{v.addrBus | aton 2}}:{{v.addrSlot | aton 2}}.{{v.addrFunc}}
                {{/if}}</span>
                </td>
                <td>{{v.source}}</td>
            </tr>
        {{/each}}
        `)(data)
    }
}