import {Controller} from './controller.js'
import {DataStoreApi} from "../api/datastoresapi.js";
import {DataStoreTableWid} from "../widget/datastore/datastoretable.js";
import {FileUploadWid} from "../widget/common/fileupload.js";
import {UploadApi} from "../api/uploadapi.js";
import {CheckboxWid} from "../widget/common/checkbox.js";
import {ConfirmActionWid} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends CheckboxWid {
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
        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new DataStoreTableWid({id: this.child('#display-table')});
        this.upload = new FileUploadWid({id: props.upload});
        this.confirm = props.confirm;

        this.upload.onsubmit(this.uuids, function (e) {
            new UploadApi({uuids: e.data.store, id: '#process'}).upload(e.form);
        });
        // register buttons's  click.
        $(this.child('#delete')).on("click", this.uuids, (e) => {
            let uuids = e.data.store.slice();
            new ConfirmActionWid({
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
                this.CheckboxWid.refresh();
            });
        });
        this.table.refresh((e) => {
            this.CheckboxWid.refresh();
        });

        this.refresh();
    }

    create(data) {
        new DataStoreApi().create(data);
    }

    refresh() {
        this.table.refresh((e) => {
            this.CheckboxWid.refresh();
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
