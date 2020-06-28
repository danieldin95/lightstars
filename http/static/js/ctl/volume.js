import {Ctl} from "./ctl.js";
import {CheckBox} from "../widget/checkbox/checkbox.js";
import VolumeTable from "../widget/volume/table.js";


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
        this.uuid = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new VolumeTable({
            id: this.child('#display-table'),
            uuid: this.uuid,
            name: this.name,
        });

        // refresh table and register refresh click.
        $(this.child('#create')).on("click", (e) => {
            console.log('create')

        });
        $(this.child('#edit')).on("click", (e) => {
            console.log('edit')

        });
        $(this.child('#remove')).on("click", (e) => {
            console.log('remove')

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
            console.log('length', $(this.child('#on-this')).length)
            // register click on this table row.

            $(this.child('#on-this')).on('click', function (e) {
                console.log($(this).attr('data'))
            });
        });
    }
}
