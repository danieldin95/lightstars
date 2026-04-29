import {Widget} from "../widget.js";
import {DiskApi} from "../../api/diskapi.js";
import {DataStoreApi} from "../../api/datastoresapi.js";
import {Location} from "../../lib/location.js";


export class DiskTableWid extends Widget {
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
            $(e.data.id).html(e.data.render(e.data.formatData(e.resp)));
            $(e.data.id).find(".js-datastore-link").off("click").on("click", function (evt) {
                evt.preventDefault();
                let name = ($(this).attr("data-store") || "").trim();
                if (!name) {
                    return;
                }
                new DataStoreApi().list((resp) => {
                    let items = (resp && resp.resp && resp.resp.items) ? resp.resp.items : [];
                    let found = items.find((it) => (it.name || "").trim() === name);
                    if (!found || !found.uuid) {
                        return;
                    }
                    let query = Location.query();
                    window.location.hash = `#/datastore/${found.uuid}${query ? "?" + query : ""}`;
                });
            });
            func({data, resp: e.resp});
        });
    }

    formatSource(source) {
        let s = (source || "").toString();
        // Normalize local datastore path to a concise display style.
        // e.g. /lightstar/datastore/01/ubuntu.iso -> datastore@01:/ubuntu.iso
        let m = s.match(/^\/lightstar\/datastore\/(\d{1,2})(\/.*)?$/);
        if (m) {
            let idx = m[1].padStart(2, "0");
            let rest = m[2] || "/";
                return `datastore@${idx}:${rest}`;
            }
            return s;
    }

    formatData(data) {
        let items = (data && data.items) ? data.items : [];
        return Object.assign({}, data, {
            items: items.map((v) => {
                let src = (v.source || "").toString();
                let m = src.match(/^\/lightstar\/datastore\/(\d{1,2})(\/.*)?$/);
                return Object.assign({}, v, {
                    sourceDisplay: this.formatSource(v.source),
                    sourceStore: m
                        ? `datastore@${m[1].padStart(2, "0")}`
                        : "",
                    sourceSuffix: m ? (m[2] || "/") : "",
                });
            }),
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
                <td>
                    {{if v.sourceStore}}
                        <a href="javascript:void(0)" class="js-datastore-link" data-store="{{v.sourceStore}}">{{v.sourceStore}}</a><span>{{":" + (v.sourceSuffix || "/")}}</span>
                    {{else}}
                        {{v.sourceDisplay || v.source}}
                    {{/if}}
                </td>
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
