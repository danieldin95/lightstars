import {Widget} from "../widget.js";
import {InterfaceApi} from "../../api/interfaceapi.js";
import {NetworkApi} from "../../api/networkapi.js";
import {Location} from "../../lib/location.js";


export class InterfaceTableWid extends Widget {
    // {
    //   id: '#xx',
    //   inst: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
    }

    loading() {
        return `<tr><td colspan="7" class="text-center">Loading...</td></tr>`;
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
            $(e.data.id).html(e.data.render(e.data.formatData(e.resp)));
            $(e.data.id).find(".js-network-link").off("click").on("click", function (evt) {
                evt.preventDefault();
                let bridge = ($(this).attr("data-bridge") || "").trim();
                if (!bridge) {
                    return;
                }
                new NetworkApi().list((resp) => {
                    let items = (resp && resp.resp && resp.resp.items) ? resp.resp.items : [];
                    let found = items.find((it) => {
                        return (it.bridge || "").trim() === bridge || (it.name || "").trim() === bridge;
                    });
                    if (!found || !found.uuid) {
                        return;
                    }
                    let query = Location.query();
                    window.location.hash = `#/network/${found.uuid}${query ? "?" + query : ""}`;
                });
            });
            func({data, resp: e.resp});
        });
    }

    formatData(data) {
        let items = (data && data.items) ? data.items : [];
        return Object.assign({}, data, {
            items: items.map((v) => {
                let sourceText = v.source == "" ? (v.network == "" ? v.hostDev : v.network) : v.source;
                return Object.assign({}, v, {
                    sourceText: sourceText || "-",
                    sourceBridge: sourceText || "",
                });
            }),
        });
    }

    render(data) {
        return this.compile(`
        {{if (items.length === 0)}}
            <tr>
                <td colspan="7" class="text-center">{{'no data to display' | i}}</td>
            </tr>
        {{/if}}
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.address}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.model}}</td>
                <td>{{v.device == "" ? '-' : v.device}}</td>
                <td>{{v.address}}</td>
                <td><span>
                {{if  v.addrType == "pci"}}
                    pci:{{v.addrBus}}:{{v.addrSlot}}.{{v.addrFunc}}
                {{/if}}</span>
                </td>
                <td>
                    {{if v.sourceBridge}}
                        <a href="javascript:void(0)" class="js-network-link" data-bridge="{{v.sourceBridge}}">{{v.sourceText}}</a>
                    {{else}}
                        {{v.sourceText}}
                    {{/if}}
                </td>
            </tr>
        {{/each}}
        `, data);
    }
}
