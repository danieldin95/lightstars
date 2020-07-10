import {CheckBox as CheckBoxCom} from "../../com/checkbox.js";


export class CheckBox {
    // {
    //   id: "#networks"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxCom({
            one: this.child('#on-one'),
            all: this.child('#on-all'),
            change: (e) => {
                this.change(e);
            },
        });
        // disable firstly.
        this.change(this.uuids);
    }

    refresh() {
        this.top.refresh();
    }

    child(id) {
        return [this.id, id].join(" ")
    }

    change(from) {
        this.uuids.store = from.store;
        if (from.store.length === 0) {
            $(this.child('#remove')).attr("disabled","disabled");
            $(this.child('#delete')).attr("disabled","disabled");
        } else {
            $(this.child('#remove')).removeAttr('disabled');
            $(this.child('#delete')).removeAttr('disabled');
        }
        if (from.store.length !== 1) {
            $(this.child('#edit')).attr("disabled","disabled");
        } else {
            $(this.child('#edit')).removeAttr('disabled');
        }
    }
}
