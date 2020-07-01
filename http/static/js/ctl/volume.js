import {Ctl} from "./ctl.js";
import {CheckBox} from "../widget/checkbox/checkbox.js";
import VolumeTable from "../widget/volume/table.js";
import {VolumeApi} from "../api/volume.js";


class CheckBoxCtl extends CheckBox {
}


export class VolumeCtl extends Ctl {
    // {
    //   id: '#network #leases',
    //   uuid: uuid of network,
    //   name: name of network,
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

        // refresh table and register refresh click.
        $(this.child('#create')).on("click", (e) => {

        });
        $(this.child('#edit')).on("click", (e) => {

        });
        $(this.child('#remove')).on("click", (e) => {

        });
        $(this.child('#refresh')).on("click", (e) => {

            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });

        this.refresh()
    }

    refresh() {
        this.table.refresh((e) => {
            this.checkbox.refresh();
            // register click on this table row.

            let _this = this
            $(this.child('#on-this')).on('click', function (e) {
                let name = $(this).attr('data-name')
                let type = $(this).attr('data-type')

                if (type === "dir") {
                    _this.table.pool = name
                    _this.refresh()
                } else {

                    _this.uuids = [name]
                    new VolumeApi({
                        pool: _this.table.pool,
                        uuids: [name],
                    }).get()
                }
            });
        });
    }
}
