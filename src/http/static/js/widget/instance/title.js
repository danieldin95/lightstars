import {FormModal} from "../form/modal.js";


export class TitleSet extends FormModal {
    //
    constructor (props) {
        super(props);

        this.title = props.data.title || "";
        this.render();
        this.loading();
    }

    render() {
        super.render();
        let title = this.view.find("input[name='title']");
        title.val(this.title);
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'setting' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form>
                    <div class="form-group">
                        <label for="cpu" class="col-form-label-sm">{{'title' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="title" value=""/>
                        </div>
                    </div>
                </form>
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button name="finish-btn" class="btn btn-outline-success btn-sm">{{'finish' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
