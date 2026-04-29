import {Controller} from './controller.js'
import {DataStoreApi} from "../api/datastores.js";
import {DataStoreTable} from "../widget/datastore/table.js";
import {FileUpload} from "../widget/common/upload.js";
import {UploadApi} from "../api/upload.js";
import {CheckBox} from "../widget/common/checkbox.js";
import {ConfirmAction} from "../widget/common/confirm.js";


class CheckBoxCtl extends CheckBox {
    change(from) {
        super.change(from);
        if (from.store.length !== 1) {
            $(this.child('#upload')).attr("disabled","disabled");
        } else {
            $(this.child('#upload')).removeAttr('disabled');
        }
    }
}


export class DataStoresCtl extends Controller {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DataStoreTable({id: this.child('#display-table')});
        this.upload = new FileUpload({id: props.upload});
        this.confirm = props.confirm;

        this.upload.onsubmit(this.uuids, function (e) {
            new UploadApi({uuids: e.data.store, id: '#process'}).upload(e.form);
        });
        // register buttons's  click.
        $(this.child('#delete')).on("click", this.uuids, (e) => {
            let uuids = e.data.store.slice();
            new ConfirmAction({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                new DataStoreApi({uuids: uuids}).delete();
            });
            $(this.confirm).modal("show");
        });

        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });

        this.refresh();
    }

    create(data) {
        new DataStoreApi().create(data);
    }

    refresh() {
        this.table.refresh((e) => {
            this.checkbox.refresh();
            // register click on this table row.
            let func = this.props.onthis;
            if (func) {
                $(this.child('#on-this')).on('click', function (e) {
                    func({uuid: $(this).attr('data')});
                });
            }
        });
    }
}
