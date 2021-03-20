import {Controller} from "./controller.js";
import {CheckBox} from "../widget/common/checkbox.js";
import VolumeTable from "../widget/volume/table.js";
import {VolumeApi} from "../api/volume.js";
import {FileUpload} from "../widget/common/upload.js";
import {UploadApi} from "../api/upload.js";


class CheckBoxCtl extends CheckBox {
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

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new VolumeTable({
            id: this.child('#display-table'),
            pool: this.pool
        });
        this.upload = new FileUpload({id: props.upload});
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
            this.refresh();
        });
        this.refresh()
    }

    current(value) {
        $(this.child("#current")).text(value);
    }

    refresh() {
        this.table.refresh((e) => {
            this.checkbox.refresh();
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
