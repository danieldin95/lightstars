import {Controller} from "./controller.js";
import {CheckboxWid} from "../widget/common/checkbox.js";
import {VolumeTableWid} from "../widget/volume/volumetable.js";
import {VolumeApi} from "../api/volumeapi.js";
import {FileUploadWid} from "../widget/common/fileupload.js";
import {UploadApi} from "../api/uploadapi.js";


class CheckBoxCtl extends CheckboxWid {
    change(from) {
        super.change(from);
    }
}


export class VolumeCtl extends Controller {
    // {
    //   id: '#pool #volume',
    //   uuid: uuid of pool,
    //   name: name of pool,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.pool = props.uuid;

        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new VolumeTableWid({
            id: this.child('#display-table'),
            pool: this.pool
        });
        this.upload = new FileUploadWid({id: props.upload});
        this.upload.onsubmit((e) => {
            new UploadApi({
                uuids: this.table.pool,
                id: '#process'
            }).upload(e.form);
        });

        // refresh table and register refresh click.
        $(this.child('#create')).on("click", (e) => {
            console.log("todo");
        });
        $(this.child('#remove')).on("click", (e) => {
            let data = {
                pool: this.table.pool,
                uuids: this.uuids.store,
            };
            if (this.props.onRemove) {
                this.props.onRemove(data);
            } else {
                new VolumeApi(data).delete();
            }
        });
        $(this.child('#refresh')).on("click", (e) => {
            this.refresh();
        });
        $(this.child("#datastore")).on("click", (e) => {
            this.table.pool = this.pool;
            this.current("");
            this.refresh();
        });
        $(this.child("#current")).on("click", (e) => {
            this.table.pool = this.pool;
            this.current("");
            this.refresh();
        });
        this.refresh()
    }

    current(value) {
        const cur = (value || "").trim();
        const curBtn = $(this.child("#current"));
        if (cur === "") {
            curBtn.text("");
            curBtn.attr("disabled", "disabled");
            return;
        }
        curBtn.text(`/${cur}`);
        curBtn.removeAttr("disabled");
    }

    refresh() {
        this.table.refresh((e) => {
            this.CheckboxWid.refresh();
            // register click on this table row.
            $(this.child('#on-this')).on('click', this, function (e) {
                let name = $(this).attr('data-name');
                let type = $(this).attr('data-type');

                if (type === "dir") {
                    e.data.table.pool = `.${name}`;
                    e.data.current(name);
                    e.data.refresh();
                } else {
                    e.data.uuids = [name];
                    new VolumeApi({
                        pool: e.data.table.pool,
                        uuids: name,
                    }).get(e.data, ()=>{});
                }
            });
        });
    }
}
